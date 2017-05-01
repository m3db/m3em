package m3db

import (
	"testing"

	"google.golang.org/grpc"

	"github.com/m3db/m3em/generated/proto/m3em"
	"github.com/m3db/m3em/node"
	mocknode "github.com/m3db/m3em/node/mocks"

	"github.com/golang/mock/gomock"
	m3dbrpc "github.com/m3db/m3db/generated/thrift/rpc"
	"github.com/stretchr/testify/require"
)

func newTestOptions() Options {
	return NewOptions(nil).
		SetNodeOptions(
			node.NewOptions(nil).
				SetOperatorClientFn(func() (*grpc.ClientConn, m3em.OperatorClient, error) { return nil, nil, nil }),
		)
}

func TestHealthEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockM3DBClient := m3dbrpc.NewMockTChanNode(ctrl)
	mockM3DBClient.EXPECT().Health(gomock.Any()).Return(&m3dbrpc.NodeHealthResult_{
		Bootstrapped: true,
		Ok:           false,
		Status:       "NOT_OK",
	}, nil)

	opts := newTestOptions()
	mockNode := mocknode.NewMockServiceNode(ctrl)

	nodeInterface, err := New(mockNode, opts)
	require.NoError(t, err)
	testNode := nodeInterface.(*m3dbNode)
	testNode.m3dbClient = mockM3DBClient

	health, err := testNode.Health()
	require.NoError(t, err)
	require.True(t, health.Bootstrapped)
	require.False(t, health.OK)
	require.Equal(t, "NOT_OK", health.Status)
}
