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

package node

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"

	"github.com/m3db/m3em/build"
	"github.com/m3db/m3em/generated/proto/m3em"
	mtime "github.com/m3db/m3em/time"

	"github.com/m3db/m3cluster/services"
	"github.com/m3db/m3cluster/shard"
	m3dbrpc "github.com/m3db/m3db/generated/thrift/rpc"
	m3dbchannel "github.com/m3db/m3db/network/server/tchannelthrift/node/channel"
	xerrors "github.com/m3db/m3x/errors"
	"github.com/m3db/m3x/log"
	gu "github.com/nu7hatch/gouuid"
	tchannel "github.com/uber/tchannel-go"
	"github.com/uber/tchannel-go/thrift"
)

var (
	errUnableToSetupInitializedNode = fmt.Errorf("unable to setup node, must be either setup/uninitialized")
	errUnableToResetNode            = fmt.Errorf("unable to reset node, must be either setup/running")
	errUnableToTeardownNode         = fmt.Errorf("unable to teardown node, must be either setup/running")
	errUnableToStartNode            = fmt.Errorf("unable to start node, it must be setup")
	errUnableToStopNode             = fmt.Errorf("unable to stop node, it must be running")
	errUnableToOverrideConf         = fmt.Errorf("unable to override node configuration, it must be setup")
)

type svcNode struct {
	sync.Mutex
	logger            xlog.Logger
	opts              NodeOptions
	id                string
	rack              string
	zone              string
	weight            uint32
	endpoint          string
	shards            shard.Shards
	status            NodeStatus
	currentBuild      build.ServiceBuild
	currentConf       build.ServiceConfiguration
	clientConn        *grpc.ClientConn
	client            m3em.OperatorClient
	listeners         *listenerGroup
	heartbeater       *opHeartbeatServer
	operatorUUID      string
	heartbeatEndpoint string

	m3dbClient m3dbrpc.TChanNode
}

// NewServiceNode returns a new ServiceNode.
func NewServiceNode(
	node services.PlacementInstance,
	opts NodeOptions,
) (ServiceNode, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	clientConn, client, err := opts.OperatorClientFn()()
	if err != nil {
		return nil, err
	}

	uuid, err := gu.NewV4()
	if err != nil {
		return nil, err
	}

	var (
		listeners      = newListenerGroup()
		hbUUID         = string(uuid[:])
		heartbeater    *opHeartbeatServer
		routerEndpoint string
	)

	if opts.HeartbeatOptions().Enabled() {
		router := opts.HeartbeatOptions().HeartbeatRouter()
		routerEndpoint = router.Endpoint()
		heartbeater = newHeartbeater(listeners, opts.HeartbeatOptions(), opts.InstrumentOptions())
		if err := router.Register(hbUUID, heartbeater); err != nil {
			return nil, fmt.Errorf("unable to register heartbeat server with router: %v", err)
		}
	}

	retNode := &svcNode{
		logger:            opts.InstrumentOptions().Logger(),
		opts:              opts,
		id:                node.ID(),
		rack:              node.Rack(),
		zone:              node.Zone(),
		weight:            node.Weight(),
		endpoint:          node.Endpoint(),
		shards:            node.Shards(),
		status:            NodeStatusUninitialized,
		listeners:         listeners,
		client:            client,
		clientConn:        clientConn,
		heartbeater:       heartbeater,
		heartbeatEndpoint: routerEndpoint,
		operatorUUID:      hbUUID,
	}
	return retNode, nil
}

func (i *svcNode) String() string {
	i.Lock()
	defer i.Unlock()
	return fmt.Sprintf(
		"ServiceNode[ID=%s, Rack=%s, Zone=%s, Weight=%d, Endpoint=%s, Shards=%s]",
		i.id, i.rack, i.zone, i.weight, i.endpoint, i.shards.String(),
	)
}

func (i *svcNode) ID() string {
	i.Lock()
	defer i.Unlock()
	return i.id
}

func (i *svcNode) SetID(id string) services.PlacementInstance {
	i.Lock()
	defer i.Unlock()
	i.id = id
	return i
}

func (i *svcNode) Rack() string {
	i.Lock()
	defer i.Unlock()
	return i.rack
}

func (i *svcNode) SetRack(r string) services.PlacementInstance {
	i.Lock()
	defer i.Unlock()
	i.rack = r
	return i
}

func (i *svcNode) Zone() string {
	i.Lock()
	defer i.Unlock()
	return i.zone
}

func (i *svcNode) SetZone(z string) services.PlacementInstance {
	i.Lock()
	defer i.Unlock()
	i.zone = z
	return i
}

func (i *svcNode) Weight() uint32 {
	i.Lock()
	defer i.Unlock()
	return i.weight
}

func (i *svcNode) SetWeight(w uint32) services.PlacementInstance {
	i.Lock()
	defer i.Unlock()
	i.weight = w
	return i
}

func (i *svcNode) Endpoint() string {
	i.Lock()
	defer i.Unlock()
	return i.endpoint
}

func (i *svcNode) SetEndpoint(ip string) services.PlacementInstance {
	i.Lock()
	defer i.Unlock()
	i.endpoint = ip
	return i
}

func (i *svcNode) Shards() shard.Shards {
	i.Lock()
	defer i.Unlock()
	return i.shards
}

func (i *svcNode) SetShards(s shard.Shards) services.PlacementInstance {
	i.Lock()
	defer i.Unlock()
	i.shards = s
	return i
}

func (i *svcNode) Setup(
	bld build.ServiceBuild,
	conf build.ServiceConfiguration,
	token string,
	force bool,
) error {
	i.Lock()
	defer i.Unlock()
	if i.status != NodeStatusUninitialized &&
		i.status != NodeStatusSetup {
		return errUnableToSetupInitializedNode
	}

	i.currentConf = conf
	i.currentBuild = bld

	freq := uint32(i.opts.HeartbeatOptions().Interval().Seconds())
	err := i.opts.Retrier().Attempt(func() error {
		ctx := context.Background()
		_, err := i.client.Setup(ctx, &m3em.SetupRequest{
			OperatorUuid:           i.operatorUUID,
			SessionToken:           token,
			Force:                  force,
			HeartbeatEnabled:       i.opts.HeartbeatOptions().Enabled(),
			HeartbeatEndpoint:      i.heartbeatEndpoint,
			HeartbeatFrequencySecs: freq,
		})
		return err
	})

	if err != nil {
		return fmt.Errorf("unable to setup: %v", err)
	}

	// TODO(prateek): make heartbeat pickup existing agent state

	// Wait till we receive our first heartbeat
	if i.opts.HeartbeatOptions().Enabled() {
		i.logger.Infof("waiting until initial heartbeat is recieved")
		received := mtime.WaitUntil(i.heartbeatReceived, i.opts.HeartbeatOptions().Timeout())
		if !received {
			return fmt.Errorf("did not receive heartbeat response from remote agent within timeout")
		}
		i.logger.Infof("initial heartbeat recieved")

		// start hb monitoring
		if err := i.heartbeater.start(); err != nil {
			return fmt.Errorf("unable to start heartbeat monitor loop: %v", err)
		}
	}

	// transfer build
	if err := i.opts.Retrier().Attempt(func() error {
		return i.sendFile(bld, m3em.FileType_SERVICE_BINARY, force)
	}); err != nil {
		return fmt.Errorf("unable to transfer build: %v", err)
	}

	if err := i.opts.Retrier().Attempt(func() error {
		return i.sendFile(conf, m3em.FileType_SERVICE_CONFIG, force)
	}); err != nil {
		return fmt.Errorf("unable to transfer config: %v", err)
	}

	i.status = NodeStatusSetup
	return nil
}

func (i *svcNode) heartbeatReceived() bool {
	return !i.heartbeater.lastHeartbeatTime().IsZero()
}

func (i *svcNode) sendFile(
	file build.IterableBytesWithID,
	fileType m3em.FileType,
	overwrite bool,
) error {
	filename := file.ID()
	iter, err := file.Iter(i.opts.TransferBufferSize())
	if err != nil {
		return err
	}
	defer iter.Close()

	ctx := context.Background()
	stream, err := i.client.Transfer(ctx)
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

func (i *svcNode) Teardown() error {
	i.Lock()
	defer i.Unlock()
	if status := i.status; status != NodeStatusRunning &&
		status != NodeStatusSetup &&
		status != NodeStatusError {
		return errUnableToTeardownNode
	}

	// clear any listeners
	i.listeners.clear()

	if err := i.opts.Retrier().Attempt(func() error {
		ctx := context.Background()
		_, err := i.client.Teardown(ctx, &m3em.TeardownRequest{})
		return err
	}); err != nil {
		return err
	}

	if err := i.Close(); err != nil {
		return err
	}

	i.status = NodeStatusUninitialized
	return nil
}

func (i *svcNode) Close() error {
	var err xerrors.MultiError

	if conn := i.clientConn; conn != nil {
		i.clientConn = nil
		err = err.Add(conn.Close())
	}

	if hbServer := i.heartbeater; hbServer != nil {
		hbServer.stop()
		err = err.Add(i.opts.HeartbeatOptions().HeartbeatRouter().Deregister(i.operatorUUID))
		i.heartbeater = nil
		i.operatorUUID = ""
	}

	return err.FinalError()
}

func (i *svcNode) Start() error {
	i.Lock()
	defer i.Unlock()
	if i.status != NodeStatusSetup {
		return errUnableToStartNode
	}

	if err := i.opts.Retrier().Attempt(func() error {
		ctx := context.Background()
		_, err := i.client.Start(ctx, &m3em.StartRequest{})
		return err
	}); err != nil {
		return err
	}

	i.status = NodeStatusRunning
	return nil
}

func (i *svcNode) Stop() error {
	i.Lock()
	defer i.Unlock()
	if i.status != NodeStatusRunning {
		return errUnableToStopNode
	}

	if err := i.opts.Retrier().Attempt(func() error {
		ctx := context.Background()
		_, err := i.client.Stop(ctx, &m3em.StopRequest{})
		return err
	}); err != nil {
		return err
	}

	i.status = NodeStatusSetup
	return nil
}

func (i *svcNode) Status() NodeStatus {
	i.Lock()
	defer i.Unlock()
	return i.status
}

func (i *svcNode) RegisterListener(l Listener) ListenerID {
	return ListenerID(i.listeners.add(l))
}

func (i *svcNode) DeregisterListener(token ListenerID) {
	i.listeners.remove(int(token))
}

func (i *svcNode) thriftClient() (m3dbrpc.TChanNode, error) {
	i.Lock()
	defer i.Unlock()
	if i.m3dbClient != nil {
		return i.m3dbClient, nil
	}
	channel, err := tchannel.NewChannel("Client", nil)
	if err != nil {
		return nil, fmt.Errorf("could not create new tchannel channel: %v", err)
	}
	endpoint := &thrift.ClientOptions{HostPort: i.endpoint}
	thriftClient := thrift.NewClient(channel, m3dbchannel.ChannelName, endpoint)
	client := m3dbrpc.NewTChanNodeClient(thriftClient)
	i.m3dbClient = client
	return i.m3dbClient, nil
}

func (i *svcNode) Health() (ServiceNodeHealth, error) {
	healthResult := ServiceNodeHealth{}

	client, err := i.thriftClient()
	if err != nil {
		return healthResult, err
	}

	attemptFn := func() error {
		tctx, _ := thrift.NewContext(i.opts.OperationTimeout())
		result, err := client.Health(tctx)
		if err != nil {
			return err
		}
		healthResult.Bootstrapped = result.GetBootstrapped()
		healthResult.OK = result.GetOk()
		healthResult.Status = result.GetStatus()
		return nil
	}

	retrier := i.opts.Retrier()
	err = retrier.Attempt(attemptFn)
	return healthResult, err
}
