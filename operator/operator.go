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

package operator

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/m3db/m3em/build"
	"github.com/m3db/m3em/generated/proto/m3em"
	"github.com/m3db/m3em/os/fs"

	"github.com/m3db/m3x/instrument"
	"github.com/m3db/m3x/log"
	"google.golang.org/grpc"
)

type operator struct {
	endpoint    string
	clientConn  *grpc.ClientConn
	client      m3em.OperatorClient
	opts        Options
	logger      xlog.Logger
	listeners   *listenerGroup
	heartbeater *opHeartbeater
}

// New creates a new operator
func New(
	endpoint string,
	opts Options,
) (Operator, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(endpoint, grpc.WithTimeout(opts.Timeout()), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	var (
		client      = m3em.NewOperatorClient(conn)
		listeners   = newListenerGroup()
		heartbeater = newHeartbeater(client, listeners, opts.InstrumentOptions(), opts.HeartbeatOptions())
	)

	return &operator{
		endpoint:    endpoint,
		listeners:   listeners,
		client:      client,
		clientConn:  conn,
		opts:        opts,
		logger:      opts.InstrumentOptions().Logger(),
		heartbeater: heartbeater,
	}, nil
}

func (o *operator) Setup(
	bld build.ServiceBuild,
	conf build.ServiceConfiguration,
	token string,
	force bool,
) error {
	if err := o.sendSetup(token, force); err != nil {
		return fmt.Errorf("unable to setup: %v", err)
	}

	if o.opts.HeartbeatOptions().Enabled() {
		if err := o.heartbeater.start(); err != nil {
			return fmt.Errorf("unable to heartbeat: %v", err)
		}
	}

	if err := o.sendBuild(bld, force); err != nil {
		return fmt.Errorf("unable to transfer build: %v", err)
	}

	if err := o.sendConfig(conf, force); err != nil {
		return fmt.Errorf("unable to transfer config: %v", err)
	}
	return nil
}

func (o *operator) sendSetup(token string, force bool) error {
	return o.opts.Retrier().Attempt(func() error {
		ctx := context.Background()
		_, err := o.client.Setup(ctx, &m3em.SetupRequest{
			Token: token,
			Force: force,
		})
		return err
	})
}

func (o *operator) sendBuild(
	bld build.ServiceBuild,
	overwrite bool,
) error {
	return o.opts.Retrier().Attempt(func() error {
		filename := bld.SourcePath()
		iter, err := fs.NewSizedFileReaderIter(filename, o.opts.TransferBufferSize())
		if err != nil {
			return err
		}
		defer iter.Close()
		return o.sendFile(iter, m3em.FileType_M3DB_BINARY, bld.ID(), overwrite)
	})
}

func (o *operator) sendConfig(
	conf build.ServiceConfiguration,
	overwrite bool,
) error {
	return o.opts.Retrier().Attempt(func() error {
		bytes, err := conf.MarshalText()
		if err != nil {
			return err
		}
		iter := fs.NewBytesReaderIter(bytes)
		return o.sendFile(iter, m3em.FileType_M3DB_CONFIG, conf.ID(), overwrite)
	})
}

func (o *operator) sendFile(
	iter fs.FileReaderIter,
	fileType m3em.FileType,
	filename string,
	overwrite bool,
) error {
	ctx := context.Background()
	stream, err := o.client.Transfer(ctx)
	if err != nil {
		return err
	}
	chunkIdx := 0
	for ; iter.Next(); chunkIdx++ {
		bytes := iter.Current()
		request := &m3em.TransferRequest{
			Type:       fileType,
			Filename:   filename,
			Overwrite:  overwrite,
			ChunkBytes: bytes,
			ChunkIdx:   int32(chunkIdx),
		}
		err := stream.Send(request)
		if err != nil {
			stream.CloseSend()
			return err
		}
	}
	if err := iter.Err(); err != nil {
		stream.CloseSend()
		return err
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	if int(response.NumChunksRecvd) != chunkIdx {
		return fmt.Errorf("sent %d chunks, server only received %d of them", chunkIdx, response.NumChunksRecvd)
	}

	if iter.Checksum() != response.FileChecksum {
		return fmt.Errorf("expected file checksum: %d, received: %d", iter.Checksum(), response.FileChecksum)
	}

	return nil
}

func (o *operator) Start() error {
	return o.opts.Retrier().Attempt(func() error {
		ctx := context.Background()
		_, err := o.client.Start(ctx, &m3em.StartRequest{})
		return err
	})
}

func (o *operator) Stop() error {
	return o.opts.Retrier().Attempt(func() error {
		ctx := context.Background()
		_, err := o.client.Stop(ctx, &m3em.StopRequest{})
		return err
	})
}

func (o *operator) Teardown() error {
	return o.opts.Retrier().Attempt(func() error {
		o.heartbeater.stop()
		ctx := context.Background()
		_, err := o.client.Teardown(ctx, &m3em.TeardownRequest{})
		return err
	})
}

func (o *operator) Reset() error {
	panic("not implemented")
}

func (o *operator) close() error {
	if conn := o.clientConn; conn != nil {
		o.clientConn = nil
		return conn.Close()
	}
	return nil
}

func (o *operator) RegisterListener(l Listener) ListenerID {
	return ListenerID(o.listeners.add(l))
}

func (o *operator) DeregisterListener(token ListenerID) {
	o.listeners.remove(int(token))
}

// NewListener creates a new listener
func NewListener(
	onProcessTerminate func(string),
	onHeartbeatTimeout func(last time.Time),
	onOverwrite func(string),
	onInternalError func(err error),
) Listener {
	return &listener{
		onProcessTerminate: onProcessTerminate,
		onHeartbeatTimeout: onHeartbeatTimeout,
		onOverwrite:        onOverwrite,
		onInternalError:    onInternalError,
	}
}

type listener struct {
	onProcessTerminate func(string)
	onHeartbeatTimeout func(last time.Time)
	onOverwrite        func(string)
	onInternalError    func(err error)
}

func (l *listener) OnProcessTerminate(desc string) {
	if l.onProcessTerminate != nil {
		l.onProcessTerminate(desc)
	}
}

func (l *listener) OnHeartbeatTimeout(lastHeartbeatTs time.Time) {
	if l.onHeartbeatTimeout != nil {
		l.onHeartbeatTimeout(lastHeartbeatTs)
	}
}

func (l *listener) OnOverwrite(desc string) {
	if l.onOverwrite != nil {
		l.onOverwrite(desc)
	}
}

func (l *listener) OnInternalError(err error) {
	if l.onInternalError != nil {
		l.onInternalError(err)
	}
}

type listenerGroup struct {
	sync.Mutex
	elems map[int]Listener
	token int
}

func newListenerGroup() *listenerGroup {
	return &listenerGroup{
		elems: make(map[int]Listener),
	}
}

func (lg *listenerGroup) add(l Listener) int {
	lg.Lock()
	defer lg.Unlock()
	lg.token++
	lg.elems[lg.token] = l
	return lg.token
}

func (lg *listenerGroup) remove(t int) {
	lg.Lock()
	defer lg.Unlock()
	delete(lg.elems, t)
}

func (lg *listenerGroup) notifyTimeout(lastTs time.Time) {
	lg.Lock()
	defer lg.Unlock()
	for _, l := range lg.elems {
		go l.OnHeartbeatTimeout(lastTs)
	}
}

func (lg *listenerGroup) notifyTermination(desc string) {
	lg.Lock()
	defer lg.Unlock()
	for _, l := range lg.elems {
		go l.OnProcessTerminate(desc)
	}
}

func (lg *listenerGroup) notifyOverwrite(desc string) {
	lg.Lock()
	defer lg.Unlock()
	for _, l := range lg.elems {
		go l.OnOverwrite(desc)
	}
}

func (lg *listenerGroup) notifyInternalError(err error) {
	lg.Lock()
	defer lg.Unlock()
	for _, l := range lg.elems {
		go l.OnInternalError(err)
	}
}

// Heartbeating from agent -> operator: this is to ensure capture of asynchronous
// error conditions, e.g. the child process kicked off by the agent dies, but the
// operator is not informed.
type opHeartbeater struct {
	sync.RWMutex
	wg            sync.WaitGroup
	opts          HeartbeatOptions
	iopts         instrument.Options
	listeners     *listenerGroup
	opClient      m3em.OperatorClient
	client        m3em.Operator_HeartbeatClient
	clientCancel  context.CancelFunc
	running       bool
	stopChan      chan bool
	lastHeartbeat heartbeat
}

type heartbeat struct {
	ts             time.Time
	processRunning bool
}

func newHeartbeater(
	client m3em.OperatorClient,
	lg *listenerGroup,
	iopts instrument.Options,
	opts HeartbeatOptions,
) *opHeartbeater {
	return &opHeartbeater{
		opClient:  client,
		opts:      opts,
		iopts:     iopts,
		listeners: lg,
		stopChan:  make(chan bool, 1),
	}
}

func (h *opHeartbeater) start() error {
	h.Lock()
	defer h.Unlock()
	if h.running {
		return fmt.Errorf("already heartbeating, terminate existing process first")
	}

	var (
		ctx, cancel = context.WithCancel(context.Background())
		request     = &m3em.HeartbeatRequest{
			FrequencySecs: uint32(h.opts.Interval().Seconds()),
		}
	)
	heartbeatClient, err := h.opClient.Heartbeat(ctx, request)
	if err != nil {
		cancel()
		return err
	}

	h.running = true
	h.client = heartbeatClient
	h.clientCancel = cancel
	h.wg.Add(2)
	go h.heartbeatLoop()
	go h.monitorLoop()
	return nil
}

func (h *opHeartbeater) lastHeartbeatTime() time.Time {
	h.RLock()
	defer h.RUnlock()
	return h.lastHeartbeat.ts
}

func (h *opHeartbeater) monitorLoop() {
	defer h.wg.Done()
	var (
		monitor         = true
		checkInterval   = h.opts.CheckInterval()
		timeoutInterval = h.opts.Timeout()
		nowFn           = h.opts.NowFn()
		lastCheck       time.Time
	)

	for monitor {
		select {
		case <-h.stopChan:
			monitor = false
		default:
			last := h.lastHeartbeatTime()
			if last == lastCheck {
				time.Sleep(checkInterval)
				continue
			}
			lastCheck = last
			if !lastCheck.IsZero() && nowFn().Sub(lastCheck) > timeoutInterval {
				h.listeners.notifyTimeout(lastCheck)
			}
			time.Sleep(checkInterval)
		}
	}
}

func (h *opHeartbeater) heartbeatLoop() {
	defer h.wg.Done()
	for {
		msg, err := h.client.Recv()
		if err == context.Canceled {
			h.iopts.Logger().Infof("heartbeating cancelled")
			break
		} else if err == io.EOF || err != nil {
			heartbeatErr := fmt.Errorf("unexpected error while streaming heartbeat msgs: %v", err)
			h.iopts.Logger().Warn(heartbeatErr.Error())
			h.listeners.notifyInternalError(heartbeatErr)
			break
		}
		h.processHeartbeatMsg(msg)
	}
}

func (h *opHeartbeater) updateLastHeartbeat(hb heartbeat) {
	h.Lock()
	h.lastHeartbeat = hb
	h.Unlock()
}

func (h *opHeartbeater) processHeartbeatMsg(msg *m3em.HeartbeatResponse) {
	nowFn := h.opts.NowFn()

	switch msg.GetCode() {
	case m3em.HeartbeatCode_HEARTBEAT_CODE_HEALTHY:
		newHeartbeat := heartbeat{
			ts:             nowFn(),
			processRunning: msg.GetProcessRunning(),
		}
		h.updateLastHeartbeat(newHeartbeat)

	case m3em.HeartbeatCode_HEARTBEAT_CODE_PROCESS_TERMINATION:
		newHeartbeat := heartbeat{
			ts:             nowFn(),
			processRunning: false,
		}
		h.updateLastHeartbeat(newHeartbeat)
		h.listeners.notifyTermination(msg.GetError())

	case m3em.HeartbeatCode_HEARTBEAT_CODE_OVERWRITTEN:
		h.listeners.notifyOverwrite(msg.GetError())
		h.stop()

	default:
		err := fmt.Errorf("received unknown heartbeat msg: %v", *msg)
		h.listeners.notifyInternalError(err)
		h.stop()
	}
}

func (h *opHeartbeater) stop() {
	h.Lock()
	if !h.running {
		h.resetWithLock()
		h.Unlock()
		return
	}

	h.clientCancel()
	h.stopChan <- true
	h.Unlock()
	h.wg.Wait()

	h.Lock()
	h.resetWithLock()
	h.Unlock()
}

func (h *opHeartbeater) resetWithLock() {
	h.running = false
	h.client = nil
	h.clientCancel = nil
}
