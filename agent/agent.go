// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package agent

import (
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"

	hb "github.com/m3db/m3em/generated/proto/heartbeat"
	"github.com/m3db/m3em/generated/proto/m3em"
	"github.com/m3db/m3em/os/exec"
	"github.com/m3db/m3em/os/fs"

	xclock "github.com/m3db/m3x/clock"
	xerrors "github.com/m3db/m3x/errors"
	"github.com/m3db/m3x/instrument"
	xlog "github.com/m3db/m3x/log"
	"github.com/uber-go/tally"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	defaultReportInterval       = 5 * time.Second
	defaultHeartbeatConnTimeout = time.Minute
	defaultTestCanaryPrefix     = "test-canary-file"
)

var (
	errProcessMonitorNotDefined = fmt.Errorf("process monitor not defined")
)

type opAgent struct {
	sync.RWMutex
	token               string
	executablePath      string
	configPath          string
	newProcessMonitorFn newProcessMonitorFn
	processMonitor      exec.ProcessMonitor
	heartbeater         *opAgentHeartBeater

	running            int32
	stopping           int32
	heartbeatTimeoutCh chan struct{}

	opts    Options
	logger  xlog.Logger
	metrics *opAgentMetrics
	doneCh  chan struct{}
	closeCh chan struct{}
}

type newProcessMonitorFn func(exec.Cmd, exec.ProcessListener) (exec.ProcessMonitor, error)

// New creates and retuns a new Operator Agent
func New(
	opts Options,
) (Agent, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	if err := canaryWriteTest(opts.WorkingDirectory()); err != nil {
		return nil, err
	}

	agent := &opAgent{
		opts:                opts,
		logger:              opts.InstrumentOptions().Logger(),
		metrics:             newAgentMetrics(opts.InstrumentOptions().MetricsScope()),
		newProcessMonitorFn: exec.NewProcessMonitor,
		doneCh:              make(chan struct{}),
		closeCh:             make(chan struct{}),
	}
	go agent.reportMetrics()
	return agent, nil
}

func (o *opAgent) Close() error {
	o.closeCh <- struct{}{}
	<-o.doneCh
	return nil
}

func canaryWriteTest(dir string) error {
	fi, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("unable to stat directory, [ err = %v ]", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	fd, err := ioutil.TempFile(dir, defaultTestCanaryPrefix)
	if err != nil {
		return fmt.Errorf("unable to create canary file, [ err = %v ]", err)
	}
	os.Remove(fd.Name())

	return nil
}

func updateBoolGauge(b bool, m tally.Gauge) {
	if b {
		m.Update(1)
	} else {
		m.Update(0)
	}
}

func (o *opAgent) reportMetrics() {
	reportTicker := time.NewTicker(defaultReportInterval)
	for {
		select {
		case <-reportTicker.C:
			running, exec, conf := o.state()
			updateBoolGauge(running, o.metrics.running)
			updateBoolGauge(exec != "", o.metrics.execTransferred)
			updateBoolGauge(conf != "", o.metrics.confTransferred)
		case <-o.closeCh:
			reportTicker.Stop()
			o.doneCh <- struct{}{}
			return
		}
	}
}

func (o *opAgent) Running() bool {
	return atomic.LoadInt32(&o.running) == 1
}

func (o *opAgent) state() (bool, string, string) {
	o.RLock()
	defer o.RUnlock()
	return o.Running(), o.executablePath, o.configPath
}

func (o *opAgent) initFile(
	fileType m3em.FileType,
	filename string,
	overwrite bool,
) (*os.File, error) {
	targetFile := path.Join(o.opts.WorkingDirectory(), filename)
	if _, err := os.Stat(targetFile); os.IsExist(err) && !overwrite {
		return nil, err
	}

	flags := os.O_CREATE | os.O_WRONLY
	if overwrite {
		flags = flags | os.O_TRUNC
	}

	perms := os.FileMode(0644)
	if fileType == m3em.FileType_SERVICE_BINARY {
		perms = os.FileMode(0755)
	}
	return os.OpenFile(targetFile, flags, perms)
}

func (o *opAgent) markFileDone(
	fileType m3em.FileType,
	filename string,
) {
	o.Lock()
	defer o.Unlock()
	switch fileType {
	case m3em.FileType_SERVICE_BINARY:
		o.executablePath = filename
	case m3em.FileType_SERVICE_CONFIG:
		o.configPath = filename
	default:
		o.logger.Warnf("received unknown fileType: %v, filename: %s", fileType, filename)
	}
}

func (o *opAgent) Start(ctx context.Context, request *m3em.StartRequest) (*m3em.StartResponse, error) {
	o.logger.Infof("received Start()")
	o.Lock()
	defer o.Unlock()

	if o.Running() {
		return nil, grpc.Errorf(codes.FailedPrecondition, "already running")
	}

	if o.executablePath == "" {
		return nil, grpc.Errorf(codes.FailedPrecondition, "agent missing build")
	}

	if o.configPath == "" {
		return nil, grpc.Errorf(codes.FailedPrecondition, "agent missing config")
	}

	if err := o.startWithLock(); err != nil {
		return nil, grpc.Errorf(codes.Internal, "unable to start: %v", err)
	}

	return &m3em.StartResponse{}, nil
}

func (o *opAgent) newListenerFn() func(error) {
	return func(err error) {
		if err == nil {
			err = fmt.Errorf("test process terminated without error")
		} else {
			err = fmt.Errorf("test process terminated with error: %v", err)
		}
		o.logger.Warnf("%v", err)
		if stopping := atomic.LoadInt32(&o.stopping); stopping == 0 && o.heartbeater != nil {
			o.heartbeater.notifyProcessTermination(err)
		}
		atomic.StoreInt32(&o.running, 0)
	}
}

func (o *opAgent) newProcessListener() exec.ProcessListener {
	var (
		listenFn   = o.newListenerFn()
		onComplete = func() { listenFn(nil) }
		onError    = func(err error) { listenFn(err) }
	)
	return exec.NewProcessListener(onComplete, onError)
}

func (o *opAgent) startWithLock() error {
	var (
		path, args = o.opts.ExecGenFn()(o.executablePath, o.configPath)
		osArgs     = append([]string{path}, args...)
		cmd        = exec.Cmd{
			Path:      path,
			Args:      osArgs,
			OutputDir: o.opts.WorkingDirectory(),
			Env:       o.opts.EnvMap(),
		}
		listener = o.newProcessListener()
	)
	pm, err := o.newProcessMonitorFn(cmd, listener)
	if err != nil {
		return err
	}
	o.logger.Infof("executing command: %+v", cmd)
	if err := pm.Start(); err != nil {
		return err
	}
	atomic.StoreInt32(&o.running, 1)
	o.processMonitor = pm
	return nil
}

func (o *opAgent) Stop(ctx context.Context, request *m3em.StopRequest) (*m3em.StopResponse, error) {
	o.logger.Infof("received Stop()")
	o.Lock()
	defer o.Unlock()

	if !o.Running() {
		return nil, grpc.Errorf(codes.FailedPrecondition, "not running")
	}

	atomic.StoreInt32(&o.stopping, 1)
	if err := o.stopWithLock(); err != nil {
		return nil, grpc.Errorf(codes.Internal, "unable to stop: %v", err)
	}
	atomic.StoreInt32(&o.stopping, 0)

	return &m3em.StopResponse{}, nil
}

func (o *opAgent) stopWithLock() error {
	if o.processMonitor == nil {
		return errProcessMonitorNotDefined
	}

	if err := o.processMonitor.Stop(); err != nil {
		return err
	}

	o.processMonitor = nil
	atomic.StoreInt32(&o.running, 0)
	return nil
}

func (o *opAgent) resetWithLock() error {
	var multiErr xerrors.MultiError

	if o.heartbeater != nil {
		o.logger.Infof("stopping heartbeating")
		multiErr = multiErr.Add(o.heartbeater.close())
		o.heartbeater = nil
	}

	if o.heartbeatTimeoutCh != nil {
		close(o.heartbeatTimeoutCh)
		o.heartbeatTimeoutCh = nil
	}

	if o.Running() {
		o.logger.Infof("process running, stopping")
		if err := o.stopWithLock(); err != nil {
			o.logger.Warnf("unable to stop: %v", err)
			multiErr = multiErr.Add(err)
		}
	}

	o.logger.Infof("releasing host resources")
	if err := o.opts.ReleaseHostResourcesFn()(); err != nil {
		o.logger.Infof("unable to release host resources: %v", err)
		multiErr = multiErr.Add(err)
	}

	o.token = ""
	o.executablePath = ""
	o.configPath = ""
	atomic.StoreInt32(&o.running, 0)

	return multiErr.FinalError()
}

func (o *opAgent) Teardown(ctx context.Context, request *m3em.TeardownRequest) (*m3em.TeardownResponse, error) {
	o.logger.Infof("received Teardown()")
	o.Lock()
	defer o.Unlock()

	if err := o.resetWithLock(); err != nil {
		return nil, grpc.Errorf(codes.Internal, "unable to teardown: %v", err)
	}

	return &m3em.TeardownResponse{}, nil
}

func (o *opAgent) isSetup() bool {
	o.RLock()
	defer o.RUnlock()
	return o.isSetupWithLock()
}

func (o *opAgent) isSetupWithLock() bool {
	return o.token != ""
}

func (o *opAgent) Setup(ctx context.Context, request *m3em.SetupRequest) (*m3em.SetupResponse, error) {
	o.logger.Infof("receieved Setup()")

	// nil check
	if request == nil || request.SessionToken == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "nil request")
	}

	o.Lock()
	defer o.Unlock()

	if o.token != "" && o.token != request.SessionToken && !request.Force {
		return nil, grpc.Errorf(codes.AlreadyExists, "agent already initialized with token: %s", o.token)
	}

	// stop existing heartbeating if any
	if o.heartbeater != nil {
		o.heartbeater.notifyOverwrite(fmt.Errorf("heartbeating being overwritten by new setup request: %+v", *request))
	}

	if o.isSetupWithLock() {
		// reset agent
		if err := o.resetWithLock(); err != nil {
			return nil, grpc.Errorf(codes.Aborted, "unable to reset: %v", err)
		}
	}

	// remove any files stored in the working directory
	wd := o.opts.WorkingDirectory()
	o.logger.Infof("removing contents from working directory: %s", wd)
	if err := fs.RemoveContents(wd); err != nil {
		return nil, grpc.Errorf(codes.Internal, "unable to clear working directory: %v", err)
	}

	// initialize any resources needed on the host
	o.logger.Infof("initializing host resources")
	if err := o.opts.InitHostResourcesFn()(); err != nil {
		o.resetWithLock() // release any resources
		return nil, grpc.Errorf(codes.Internal, "unable to initialize host resources: %v", err)
	}

	// setup new heartbeating
	if request.HeartbeatEnabled {
		beater, err := newHeartbeater(o, request.OperatorUuid, request.HeartbeatEndpoint, o.opts.InstrumentOptions())
		if err != nil {
			o.resetWithLock() // release any resources
			return nil, grpc.Errorf(codes.Aborted, "unable to start heartbeating process: %v", err)
		}
		o.heartbeater = beater
		o.heartbeatTimeoutCh = make(chan struct{})
		go o.monitorHeartbeatTimeout()
		o.heartbeater.start(time.Second * time.Duration(request.HeartbeatFrequencySecs))
	}

	o.token = request.SessionToken
	return &m3em.SetupResponse{}, nil
}

func (o *opAgent) monitorHeartbeatTimeout() {
	for range o.heartbeatTimeoutCh {
		o.logger.Warnf("heartbeat sending timed out, resetting agent")
		o.Lock()
		err := o.resetWithLock()
		o.Unlock()
		if err == nil {
			o.logger.Infof("successfully reset agent")
		} else {
			o.logger.Warnf("error while resetting agent: %v", err)
		}
		// resetting breaks out of the loop as it closes the timeoutCh
	}
}

func (o *opAgent) Transfer(stream m3em.Operator_TransferServer) error {
	o.logger.Infof("receieved Transfer()")
	var (
		checksum     = uint32(0)
		numChunks    = 0
		lastChunkIdx = int32(0)
		fileHandle   *os.File
		fileType     = m3em.FileType_UNKNOWN
		err          error
	)

	defer func() {
		if fileHandle != nil {
			fileHandle.Close()
		}
	}()

	for {
		request, streamErr := stream.Recv()
		if streamErr != nil && streamErr != io.EOF {
			return streamErr
		}

		if request == nil {
			o.markFileDone(fileType, fileHandle.Name())
			return stream.SendAndClose(&m3em.TransferResponse{
				FileChecksum:   checksum,
				NumChunksRecvd: int32(numChunks),
			})
		}

		if numChunks == 0 {
			// first request with any data in it, log it for visibilty
			o.logger.Infof("file transfer initiated: [ filename = %s, fileType = %s, overwrite = %s ]",
				request.GetFilename(), request.GetType().String(), request.GetOverwrite())

			fileType = request.GetType()
			fileHandle, err = o.initFile(fileType, request.GetFilename(), request.GetOverwrite())
			if err != nil {
				return err
			}
			lastChunkIdx = request.GetChunkIdx() - 1
		}

		chunkIdx := request.GetChunkIdx()
		if chunkIdx != int32(1+lastChunkIdx) {
			return fmt.Errorf("received chunkIdx: %d after %d", chunkIdx, lastChunkIdx)
		}
		lastChunkIdx = chunkIdx

		numChunks++
		bytes := request.GetChunkBytes()
		checksum = crc32.Update(checksum, crc32.IEEETable, bytes)

		numWritten, err := fileHandle.Write(bytes)
		if err != nil {
			return err
		}

		if numWritten != len(bytes) {
			return fmt.Errorf("unable to write bytes, expected: %d, observed: %d", len(bytes), numWritten)
		}

		if streamErr == io.EOF {
			o.markFileDone(fileType, fileHandle.Name())
			return stream.SendAndClose(&m3em.TransferResponse{
				FileChecksum:   checksum,
				NumChunksRecvd: int32(numChunks),
			})
		}
	}
}

type opAgentHeartBeater struct {
	sync.RWMutex
	iopts             instrument.Options
	agent             *opAgent
	running           bool
	wg                sync.WaitGroup
	msgChan           chan heartbeatMsg
	conn              *grpc.ClientConn
	client            hb.HeartbeaterClient
	operatorUUID      string
	heartbeatEndpoint string
	lastHeartbeatTs   time.Time
	nowFn             xclock.NowFn
}

func newHeartbeater(
	agent *opAgent,
	opUUID string,
	hbEndpoint string,
	iopts instrument.Options,
) (*opAgentHeartBeater, error) {
	conn, err := grpc.Dial(hbEndpoint, grpc.WithTimeout(defaultHeartbeatConnTimeout), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := hb.NewHeartbeaterClient(conn)
	return &opAgentHeartBeater{
		iopts:             iopts,
		agent:             agent,
		msgChan:           make(chan heartbeatMsg),
		conn:              conn,
		client:            client,
		operatorUUID:      opUUID,
		heartbeatEndpoint: hbEndpoint,
		nowFn:             agent.opts.NowFn(),
	}, nil
}

func (h *opAgentHeartBeater) isRunning() bool {
	h.RLock()
	defer h.RUnlock()
	return h.running
}

func (h *opAgentHeartBeater) defaultHeartbeat() hb.HeartbeatRequest {
	h.RLock()
	req := hb.HeartbeatRequest{
		OperatorUuid:   h.operatorUUID,
		ProcessRunning: h.agent.Running(),
	}
	h.RUnlock()
	return req
}

func (h *opAgentHeartBeater) start(d time.Duration) {
	h.Lock()
	h.running = true
	h.Unlock()
	h.wg.Add(1)
	go h.heartbeatLoop(d)
}

func (h *opAgentHeartBeater) sendHealthyHeartbeat() {
	beat := h.defaultHeartbeat()
	beat.Code = hb.HeartbeatCode_HEALTHY
	h.sendHeartbeat(&beat)
}

func (h *opAgentHeartBeater) heartbeatLoop(d time.Duration) {
	defer func() {
		h.Lock()
		h.running = false
		h.Unlock()
		h.wg.Done()
	}()

	// explicitly send first heartbeat as soon as we start
	h.sendHealthyHeartbeat()

	for {
		select {
		case msg := <-h.msgChan:
			if msg.stop {
				return
			}
			beat := h.defaultHeartbeat()
			msg.toRPCType(&beat)
			h.sendHeartbeat(&beat)

		default:
			h.sendHealthyHeartbeat()
			// NB(prateek): we use a sleep instead of a ticker because the latter
			// only guarantees a lower bound on the frequency of ticks, it does not
			// guarantee any upper bound. This makes it impossible to reliably detect
			// heartbeating timeouts.
			time.Sleep(d)
		}
	}
}

func (h *opAgentHeartBeater) sendHeartbeat(r *hb.HeartbeatRequest) {
	h.Lock()
	defer h.Unlock()

	logger := h.iopts.Logger()
	_, err := h.client.Heartbeat(context.Background(), r)

	// always mark the first send attempt as success to get a base time to compare against
	if err == nil || h.lastHeartbeatTs.IsZero() {
		h.lastHeartbeatTs = h.nowFn()
		return
	}

	logger.Warnf("unable to send heartbeat: %v", err)
	// check if this has been happening for past the permitted period
	timeSinceLastSend := h.nowFn().Sub(h.lastHeartbeatTs)
	timeout := h.agent.opts.HeartbeatTimeout()
	if timeSinceLastSend > timeout {
		logger.Warnf("unable to send heartbeats for %s; timing out", timeSinceLastSend.String())
		h.agent.heartbeatTimeoutCh <- struct{}{}
	}
}

func (h *opAgentHeartBeater) stop() {
	if h.isRunning() {
		h.msgChan <- heartbeatMsg{stop: true}
		h.wg.Wait()
	}
}

func (h *opAgentHeartBeater) close() error {
	h.stop()
	var err error
	if h.conn != nil {
		err = h.conn.Close()
		h.conn = nil
	}
	return err
}

func (h *opAgentHeartBeater) notifyProcessTermination(err error) {
	h.msgChan <- heartbeatMsg{
		processTerminate: true,
		err:              err,
	}
}

func (h *opAgentHeartBeater) notifyOverwrite(err error) {
	h.msgChan <- heartbeatMsg{
		overwritten: true,
		err:         err,
	}
}

type heartbeatMsg struct {
	stop             bool
	processTerminate bool
	overwritten      bool
	err              error
}

func (hm heartbeatMsg) toRPCType(req *hb.HeartbeatRequest) {
	if hm.processTerminate {
		req.Code = hb.HeartbeatCode_PROCESS_TERMINATION
		req.Error = hm.err.Error()
	} else if hm.overwritten {
		req.Code = hb.HeartbeatCode_OVERWRITTEN
		req.Error = hm.err.Error()
	} else {
		// should never happen
		req.Code = hb.HeartbeatCode_UNKNOWN
	}
}

type opAgentMetrics struct {
	// TODO(prateek): process monitor opts, metric for process uptime
	running         tally.Gauge
	execTransferred tally.Gauge
	confTransferred tally.Gauge
}

func newAgentMetrics(scope tally.Scope) *opAgentMetrics {
	subscope := scope.SubScope("agent")
	return &opAgentMetrics{
		running:         subscope.Gauge("running"),
		execTransferred: subscope.Gauge("exec_transferred"),
		confTransferred: subscope.Gauge("conf_transferred"),
	}
}
