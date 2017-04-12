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
	"sync"
	"time"

	"github.com/m3db/m3em/build"
	"github.com/m3db/m3em/generated/proto/m3em"
	"github.com/m3db/m3em/os/fs"

	"github.com/m3db/m3x/log"
	gu "github.com/nu7hatch/gouuid"
	"google.golang.org/grpc"
)

type operator struct {
	endpoint          string
	clientConn        *grpc.ClientConn
	client            m3em.OperatorClient
	opts              Options
	logger            xlog.Logger
	listeners         *listenerGroup
	heartbeater       *opHeartbeatServer
	operatorUUID      string
	heartbeatEndpoint string
}

// New creates a new operator
func New(
	agentEndpoint string,
	router HeartbeatRouter,
	opts Options,
) (Operator, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(agentEndpoint, grpc.WithTimeout(opts.Timeout()), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	uuid, err := gu.NewV4()
	if err != nil {
		return nil, err
	}

	var (
		client      = m3em.NewOperatorClient(conn)
		listeners   = newListenerGroup()
		hbUUID      = string(uuid[:])
		heartbeater *opHeartbeatServer
	)

	if opts.HeartbeatOptions().Enabled() {
		heartbeater = newHeartbeater(listeners, opts.HeartbeatOptions(), opts.InstrumentOptions())
		if err := router.Register(hbUUID, heartbeater); err != nil {
			return nil, fmt.Errorf("unable to register heartbeat server with router: %v", err)
		}
	}

	return &operator{
		endpoint:          agentEndpoint,
		listeners:         listeners,
		client:            client,
		clientConn:        conn,
		opts:              opts,
		logger:            opts.InstrumentOptions().Logger(),
		heartbeater:       heartbeater,
		heartbeatEndpoint: router.Endpoint(),
		operatorUUID:      hbUUID,
	}, nil
}

func (o *operator) Setup(
	bld build.ServiceBuild,
	conf build.ServiceConfiguration,
	token string,
	force bool,
) error {
	freq := uint32(o.opts.HeartbeatOptions().Interval().Seconds())
	err := o.opts.Retrier().Attempt(func() error {
		ctx := context.Background()
		_, err := o.client.Setup(ctx, &m3em.SetupRequest{
			OperatorUuid:           o.operatorUUID,
			SessionToken:           token,
			Force:                  force,
			HeartbeatEnabled:       o.opts.HeartbeatOptions().Enabled(),
			HeartbeatEndpoint:      o.heartbeatEndpoint,
			HeartbeatFrequencySecs: freq,
		})
		return err
	})
	if err != nil {
		return fmt.Errorf("unable to setup: %v", err)
	}

	// TODO(prateek): make heartbeat pickup existing agent state

	// Wait till we receive our first heartbeat
	if o.opts.HeartbeatOptions().Enabled() {
		o.logger.Infof("waiting until initial heartbeat is recieved")
		received := waitUntil(o.heartbeatReceived, o.opts.HeartbeatOptions().Timeout())
		if !received {
			return fmt.Errorf("did not receive heartbeat response from remote agent within timeout")
		}
		o.logger.Infof("initial heartbeat recieved")

		// start hb monitoring
		if err := o.heartbeater.start(); err != nil {
			return fmt.Errorf("unable to start heartbeat monitor loop: %v", err)
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

func (o *operator) heartbeatReceived() bool {
	return !o.heartbeater.lastHeartbeatTime().IsZero()
}

type conditionFn func() bool

func waitUntil(fn conditionFn, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if fn() {
			return true
		}
		time.Sleep(time.Second)
	}
	return false
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
) Listener {
	return &listener{
		onProcessTerminate: onProcessTerminate,
		onHeartbeatTimeout: onHeartbeatTimeout,
		onOverwrite:        onOverwrite,
	}
}

type listener struct {
	onProcessTerminate func(string)
	onHeartbeatTimeout func(last time.Time)
	onOverwrite        func(string)
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
