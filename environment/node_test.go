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

package environment

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/m3db/m3em/build"
	"github.com/m3db/m3em/generated/proto/m3em"
	mockfs "github.com/m3db/m3em/os/fs/mocks"

	"github.com/golang/mock/gomock"
	"github.com/m3db/m3cluster/services"
	m3dbrpc "github.com/m3db/m3db/generated/thrift/rpc"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const defaultRandSeed = 1234567890

var (
	defaultRandomVar = rand.New(rand.NewSource(int64(defaultRandSeed)))
)

func newMockServiceInstance(ctrl *gomock.Controller) services.PlacementInstance {
	r := defaultRandomVar
	node := services.NewMockPlacementInstance(ctrl)
	node.EXPECT().ID().AnyTimes().Return(fmt.Sprintf("%d", r.Int()))
	node.EXPECT().Endpoint().AnyTimes().Return(fmt.Sprintf("%d:%d", r.Int(), r.Int()))
	node.EXPECT().Rack().AnyTimes().Return(fmt.Sprintf("%d", r.Int()))
	node.EXPECT().Zone().AnyTimes().Return(fmt.Sprintf("%d", r.Int()))
	node.EXPECT().Weight().AnyTimes().Return(r.Uint32())
	node.EXPECT().Shards().AnyTimes().Return(nil)
	return node
}

func newMockServiceInstances(ctrl *gomock.Controller, numInstances int) []services.PlacementInstance {
	svcs := make([]services.PlacementInstance, 0, numInstances)
	for i := 0; i < numInstances; i++ {
		svcs = append(svcs, newMockServiceInstance(ctrl))
	}
	return svcs
}

func newTestNodeOptions(c *m3em.MockOperatorClient) NodeOptions {
	return NewNodeOptions(nil).
		SetOperatorClientFn(func() (*grpc.ClientConn, m3em.OperatorClient, error) {
			return nil, c, nil
		})
}

func TestNodePropertyInitialization(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	opts := newTestNodeOptions(nil)
	mockInstance := newMockServiceInstance(ctrl)
	m3dbInstance, err := NewServiceNode(mockInstance, opts)
	require.NoError(t, err)
	require.Equal(t, mockInstance.ID(), m3dbInstance.ID())
	require.Equal(t, mockInstance.Endpoint(), m3dbInstance.Endpoint())
	require.Equal(t, mockInstance.Rack(), m3dbInstance.Rack())
	require.Equal(t, mockInstance.Zone(), m3dbInstance.Zone())
	require.Equal(t, mockInstance.Weight(), m3dbInstance.Weight())
	require.Equal(t, mockInstance.Shards(), m3dbInstance.Shards())
}

func TestNodeErrorStatusIllegalTransitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestNodeOptions(mockClient)
	mockInstance := newMockServiceInstance(ctrl)
	node, err := NewServiceNode(mockInstance, opts)
	require.NoError(t, err)
	m3dbInstance := node.(*svcNode)
	require.Equal(t, NodeStatusUninitialized, m3dbInstance.Status())
	m3dbInstance.status = NodeStatusError
	require.Error(t, m3dbInstance.Start())
	require.Error(t, m3dbInstance.Stop())
	require.Error(t, m3dbInstance.Setup(nil, nil, "", false))
}

func TestNodeErrorStatusToTeardownTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestNodeOptions(mockClient)
	mockInstance := newMockServiceInstance(ctrl)
	node, err := NewServiceNode(mockInstance, opts)
	require.NoError(t, err)
	m3dbInstance := node.(*svcNode)
	require.Equal(t, NodeStatusUninitialized, m3dbInstance.Status())
	m3dbInstance.status = NodeStatusError
	mockClient.EXPECT().Teardown(gomock.Any(), gomock.Any())
	require.NoError(t, m3dbInstance.Teardown())
	require.Equal(t, NodeStatusUninitialized, m3dbInstance.Status())
}

func TestNodeUninitializedStatusIllegalTransitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestNodeOptions(mockClient)
	mockInstance := newMockServiceInstance(ctrl)
	node, err := NewServiceNode(mockInstance, opts)
	require.NoError(t, err)
	m3dbInstance := node.(*svcNode)
	require.Equal(t, NodeStatusUninitialized, m3dbInstance.Status())
	require.Error(t, m3dbInstance.Start())
	require.Error(t, m3dbInstance.Stop())
	require.Error(t, m3dbInstance.Teardown())
}

func TestNodeUninitializedStatusToSetupTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestNodeOptions(mockClient)
	mb := build.NewMockServiceBuild(ctrl)
	mc := build.NewMockServiceConfiguration(ctrl)
	mockInstance := newMockServiceInstance(ctrl)
	node, err := NewServiceNode(mockInstance, opts)
	require.NoError(t, err)
	m3dbInstance := node.(*svcNode)
	require.Equal(t, NodeStatusUninitialized, m3dbInstance.Status())

	forceSetup := false
	buildChecksum := uint32(123)
	configChecksum := uint32(321)

	dummyBytes := []byte(`some long string`)
	dummyBuildIter := mockfs.NewMockFileReaderIter(ctrl)
	gomock.InOrder(
		dummyBuildIter.EXPECT().Next().Return(true),
		dummyBuildIter.EXPECT().Current().Return(dummyBytes),
		dummyBuildIter.EXPECT().Next().Return(true),
		dummyBuildIter.EXPECT().Current().Return(dummyBytes),
		dummyBuildIter.EXPECT().Next().Return(false),
		dummyBuildIter.EXPECT().Err().Return(nil),
		dummyBuildIter.EXPECT().Checksum().Return(buildChecksum),
		dummyBuildIter.EXPECT().Close(),
	)
	mb.EXPECT().ID().Return("build-id")
	mb.EXPECT().Iter(gomock.Any()).Return(dummyBuildIter, nil)
	dummyConfIter := mockfs.NewMockFileReaderIter(ctrl)
	gomock.InOrder(
		dummyConfIter.EXPECT().Next().Return(true),
		dummyConfIter.EXPECT().Current().Return(dummyBytes),
		dummyConfIter.EXPECT().Next().Return(false),
		dummyConfIter.EXPECT().Err().Return(nil),
		dummyConfIter.EXPECT().Checksum().Return(configChecksum),
		dummyConfIter.EXPECT().Close(),
	)
	mc.EXPECT().ID().Return("config-id")
	mc.EXPECT().Iter(gomock.Any()).Return(dummyConfIter, nil)

	buildTransferClient := m3em.NewMockOperator_TransferClient(ctrl)
	gomock.InOrder(
		buildTransferClient.EXPECT().Send(&m3em.TransferRequest{
			Type:       m3em.FileType_M3DB_BINARY,
			Filename:   "build-id",
			Overwrite:  forceSetup,
			ChunkBytes: dummyBytes,
			ChunkIdx:   0,
		}).Return(nil),
		buildTransferClient.EXPECT().Send(&m3em.TransferRequest{
			Type:       m3em.FileType_M3DB_BINARY,
			Filename:   "build-id",
			Overwrite:  forceSetup,
			ChunkBytes: dummyBytes,
			ChunkIdx:   1,
		}).Return(nil),
		buildTransferClient.EXPECT().CloseAndRecv().Return(
			&m3em.TransferResponse{
				FileChecksum:   buildChecksum,
				NumChunksRecvd: 2,
			}, nil,
		),
	)
	configTransferClient := m3em.NewMockOperator_TransferClient(ctrl)
	gomock.InOrder(
		configTransferClient.EXPECT().Send(&m3em.TransferRequest{
			Type:       m3em.FileType_M3DB_CONFIG,
			Filename:   "config-id",
			Overwrite:  forceSetup,
			ChunkBytes: dummyBytes,
			ChunkIdx:   0,
		}).Return(nil),
		configTransferClient.EXPECT().CloseAndRecv().Return(
			&m3em.TransferResponse{
				FileChecksum:   configChecksum,
				NumChunksRecvd: 1,
			}, nil,
		),
	)
	gomock.InOrder(
		mockClient.EXPECT().Setup(gomock.Any(), gomock.Any()),
		mockClient.EXPECT().Transfer(gomock.Any()).Return(buildTransferClient, nil),
		mockClient.EXPECT().Transfer(gomock.Any()).Return(configTransferClient, nil),
	)

	require.NoError(t, m3dbInstance.Setup(mb, mc, "", forceSetup))
	require.Equal(t, NodeStatusSetup, m3dbInstance.Status())
	require.Equal(t, mb, m3dbInstance.currentBuild)
	require.Equal(t, mc, m3dbInstance.currentConf)
}

func TestNodeSetupStatusIllegalTransitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestNodeOptions(mockClient)
	mockInstance := newMockServiceInstance(ctrl)
	node, err := NewServiceNode(mockInstance, opts)
	require.NoError(t, err)
	m3dbInstance := node.(*svcNode)
	m3dbInstance.status = NodeStatusSetup
	require.Error(t, m3dbInstance.Stop())
}

func TestNodeSetupStatusToStartTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestNodeOptions(mockClient)
	mockInstance := newMockServiceInstance(ctrl)
	node, err := NewServiceNode(mockInstance, opts)
	require.NoError(t, err)
	m3dbInstance := node.(*svcNode)
	m3dbInstance.status = NodeStatusSetup
	mockClient.EXPECT().Start(gomock.Any(), gomock.Any())
	require.NoError(t, m3dbInstance.Start())
	require.Equal(t, NodeStatusRunning, m3dbInstance.Status())
}

func TestNodeSetupStatusToTeardownTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestNodeOptions(mockClient)
	mockInstance := newMockServiceInstance(ctrl)
	node, err := NewServiceNode(mockInstance, opts)
	require.NoError(t, err)
	m3dbInstance := node.(*svcNode)
	m3dbInstance.status = NodeStatusSetup
	mockClient.EXPECT().Teardown(gomock.Any(), gomock.Any())
	require.NoError(t, m3dbInstance.Teardown())
	require.Equal(t, NodeStatusUninitialized, m3dbInstance.Status())
}

func TestNodeRunningStatusIllegalTransitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestNodeOptions(mockClient)
	mockInstance := newMockServiceInstance(ctrl)
	node, err := NewServiceNode(mockInstance, opts)
	require.NoError(t, err)
	m3dbInstance := node.(*svcNode)
	m3dbInstance.status = NodeStatusRunning
	require.Error(t, m3dbInstance.Start())
	require.Error(t, m3dbInstance.Setup(nil, nil, "", false))
}

func TestNodeRunningStatusToStopTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestNodeOptions(mockClient)
	mockInstance := newMockServiceInstance(ctrl)
	node, err := NewServiceNode(mockInstance, opts)
	require.NoError(t, err)
	m3dbInstance := node.(*svcNode)
	m3dbInstance.status = NodeStatusRunning
	mockClient.EXPECT().Stop(gomock.Any(), gomock.Any())
	require.NoError(t, m3dbInstance.Stop())
	require.Equal(t, NodeStatusSetup, m3dbInstance.Status())
}

func TestNodeRunningStatusToTeardownTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestNodeOptions(mockClient)
	mockInstance := newMockServiceInstance(ctrl)
	node, err := NewServiceNode(mockInstance, opts)
	require.NoError(t, err)
	m3dbInstance := node.(*svcNode)
	m3dbInstance.status = NodeStatusRunning
	mockClient.EXPECT().Teardown(gomock.Any(), gomock.Any())
	require.NoError(t, m3dbInstance.Teardown())
	require.Equal(t, NodeStatusUninitialized, m3dbInstance.Status())

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

	opts := newTestNodeOptions(nil)
	mockInstance := newMockServiceInstance(ctrl)
	node, err := NewServiceNode(mockInstance, opts)
	require.NoError(t, err)
	m3dbInstance := node.(*svcNode)
	m3dbInstance.m3dbClient = mockM3DBClient

	health, err := m3dbInstance.Health()
	require.NoError(t, err)
	require.True(t, health.Bootstrapped)
	require.False(t, health.OK)
	require.Equal(t, "NOT_OK", health.Status)
}
