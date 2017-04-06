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
	"time"

	"github.com/m3db/m3em/generated/proto/m3em"
	"github.com/m3db/m3em/os/exec"

	xerrors "github.com/m3db/m3x/errors"
	xlog "github.com/m3db/m3x/log"
	"github.com/uber-go/tally"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	defaultReportInterval   = time.Duration(5 * time.Second)
	defaultTestCanaryPrefix = "test-canary-file"
)

var (
	errProcessMonitorNotDefined = fmt.Errorf("process monitor not defined")
)

type opAgent struct {
	sync.RWMutex
	opts           Options
	logger         xlog.Logger
	metrics        *opAgentMetrics
	running        bool
	token          string
	executablePath string
	configPath     string
	processMonitor exec.ProcessMonitor
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

// New creates and retuns a new Operator Agent
func New(
	opts Options,
) (Agent, error) {
	if err := canaryWriteTest(opts.WorkingDirectory()); err != nil {
		return nil, err
	}

	agent := &opAgent{
		opts:    opts,
		logger:  opts.InstrumentOptions().Logger(),
		metrics: newAgentMetrics(opts.InstrumentOptions().MetricsScope()),
	}
	go agent.reportMetrics()
	return agent, nil
}

func (o *opAgent) reportMetrics() {
	for {
		running, exec, conf := o.state()

		if running {
			o.metrics.running.Update(float64(1))
		} else {
			o.metrics.running.Update(float64(0))
		}

		if exec != "" {
			o.metrics.execTransferred.Update(float64(1))
		} else {
			o.metrics.execTransferred.Update(float64(0))
		}

		if conf != "" {
			o.metrics.confTransferred.Update(float64(1))
		} else {
			o.metrics.confTransferred.Update(float64(0))
		}

		time.Sleep(defaultReportInterval)
	}
}

func (o *opAgent) Running() bool {
	o.RLock()
	defer o.RUnlock()
	return o.running
}

func (o *opAgent) state() (bool, string, string) {
	o.RLock()
	defer o.RUnlock()
	return o.running, o.executablePath, o.configPath
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
	if fileType == m3em.FileType_M3DB_BINARY {
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
	case m3em.FileType_M3DB_BINARY:
		o.executablePath = filename
	case m3em.FileType_M3DB_CONFIG:
		o.configPath = filename
	default:
		o.logger.Warnf("received unknown fileType: %v, filename: %s", fileType, filename)
	}
}

func (o *opAgent) Heartbeat(request *m3em.HeartbeatRequest, stream m3em.Operator_HeartbeatServer) error {
	panic("not implemented")
}

func (o *opAgent) Start(ctx context.Context, request *m3em.StartRequest) (*m3em.StartResponse, error) {
	o.logger.Infof("received Start()")
	o.Lock()
	defer o.Unlock()

	if o.running {
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

func (o *opAgent) newProcessListener() exec.ProcessListener {
	onComplete := func() {
		o.Lock()
		o.logger.Warnf("test process terminated without error")
		o.running = false
		o.Unlock()
	}
	onError := func(err error) {
		o.Lock()
		o.logger.Warnf("test process terminated with error: %v", err)
		o.running = false
		o.Unlock()
	}
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
		pm, err  = exec.NewProcessMonitor(cmd, listener)
	)
	if err != nil {
		return err
	}
	o.logger.Infof("executing command: %+v", cmd)
	if err := pm.Start(); err != nil {
		return err
	}
	o.running = true
	o.processMonitor = pm
	return nil
}

func (o *opAgent) Stop(ctx context.Context, request *m3em.StopRequest) (*m3em.StopResponse, error) {
	o.logger.Infof("received Stop()")
	o.Lock()
	defer o.Unlock()

	if !o.running {
		return nil, grpc.Errorf(codes.FailedPrecondition, "not running")
	}

	if err := o.stopWithLock(); err != nil {
		return nil, grpc.Errorf(codes.Internal, "unable to stop: %v", err)
	}

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
	o.running = false
	return nil
}

func (o *opAgent) Teardown(ctx context.Context, request *m3em.TeardownRequest) (*m3em.TeardownResponse, error) {
	o.logger.Infof("received Teardown()")
	o.Lock()
	defer o.Unlock()

	var (
		multiErr xerrors.MultiError
		// workingDir = o.opts.WorkingDirectory()
		// err        = os.RemoveAll(workingDir) // TODO(prateek): does this need to be done, also, what about data directory
		// multiErr = multiErr.Add(err)
	)

	if o.processMonitor != nil {
		o.logger.Infof("processMonitor exists, attempting to Close()")
		if err := o.processMonitor.Close(); err != nil {
			o.logger.Infof("unable to Close() processMonitor: %v", err)
			multiErr = multiErr.Add(err)
		}
		o.processMonitor = nil
		o.running = false
	}

	o.logger.Infof("releasing host resources")
	if err := o.opts.ReleaseHostResourcesFn()(); err != nil {
		o.logger.Infof("unable to release host resources: %v", err)
		multiErr = multiErr.Add(err)
	}

	if finalErr := multiErr.FinalError(); finalErr != nil {
		return nil, grpc.Errorf(codes.Internal, "unable to teardown: %v", finalErr)
	}

	o.token = ""
	o.executablePath = ""
	o.configPath = ""
	o.running = false
	return &m3em.TeardownResponse{}, nil
}

func (o *opAgent) Setup(ctx context.Context, request *m3em.SetupRequest) (*m3em.SetupResponse, error) {
	if request == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "nil request")
	}
	o.Lock()
	defer o.Unlock()

	if o.token != "" && o.token != request.Token && !request.Force {
		return nil, grpc.Errorf(codes.AlreadyExists, "agent already initialized with token: %s", o.token)
	}

	if o.running {
		if err := o.stopWithLock(); err != nil {
			return nil, grpc.Errorf(codes.Aborted, "unable to stop existing process: %v", err)
		}
	}

	if o.executablePath != "" {
		if err := os.Remove(o.executablePath); err != nil {
			return nil, grpc.Errorf(codes.Aborted, "unable to remove existing executable [%s]: %v", o.executablePath, err)
		}
		o.executablePath = ""
	}

	if o.configPath != "" {
		if err := os.Remove(o.configPath); err != nil {
			return nil, grpc.Errorf(codes.Aborted, "unable to remove existing config [%s]: %v", o.configPath, err)
		}
		o.configPath = ""
	}

	if err := o.opts.InitHostResourcesFn()(); err != nil {
		return nil, grpc.Errorf(codes.Internal, "unable to initialize host resources: %v", err)
	}

	return &m3em.SetupResponse{}, nil
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

type opAgentMetrics struct {
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
