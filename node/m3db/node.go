package m3db

import (
	"fmt"
	"sync"

	"github.com/m3db/m3em/node"

	m3dbrpc "github.com/m3db/m3db/generated/thrift/rpc"
	m3dbchannel "github.com/m3db/m3db/network/server/tchannelthrift/node/channel"
	tchannel "github.com/uber/tchannel-go"
	"github.com/uber/tchannel-go/thrift"
)

type m3dbNode struct {
	sync.Mutex
	node.ServiceNode

	opts       Options
	m3dbClient m3dbrpc.TChanNode
}

// New constructs a new M3DBNode
func New(svcNode node.ServiceNode, opts Options) (Node, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return &m3dbNode{
		ServiceNode: svcNode,
		opts:        opts,
	}, nil
}

func (n *m3dbNode) thriftClient() (m3dbrpc.TChanNode, error) {
	n.Lock()
	defer n.Unlock()
	if n.m3dbClient != nil {
		return n.m3dbClient, nil
	}
	channel, err := tchannel.NewChannel("Client", nil)
	if err != nil {
		return nil, fmt.Errorf("could not create new tchannel channel: %v", err)
	}
	endpoint := &thrift.ClientOptions{HostPort: n.Endpoint()}
	thriftClient := thrift.NewClient(channel, m3dbchannel.ChannelName, endpoint)
	client := m3dbrpc.NewTChanNodeClient(thriftClient)
	n.m3dbClient = client
	return n.m3dbClient, nil
}

func (n *m3dbNode) Health() (NodeHealth, error) {
	healthResult := NodeHealth{}

	client, err := n.thriftClient()
	if err != nil {
		return healthResult, err
	}

	attemptFn := func() error {
		tctx, _ := thrift.NewContext(n.opts.NodeOptions().OperationTimeout())
		result, err := client.Health(tctx)
		if err != nil {
			return err
		}
		healthResult.Bootstrapped = result.GetBootstrapped()
		healthResult.OK = result.GetOk()
		healthResult.Status = result.GetStatus()
		return nil
	}

	retrier := n.opts.NodeOptions().Retrier()
	err = retrier.Attempt(attemptFn)
	return healthResult, err
}
