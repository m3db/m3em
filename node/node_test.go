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
	"fmt"
	"math/rand"
	"testing"

	"github.com/m3db/m3em/build"
	"github.com/m3db/m3em/generated/proto/m3em"
	mockfs "github.com/m3db/m3em/os/fs/mocks"

	"github.com/golang/mock/gomock"
	"github.com/m3db/m3cluster/services"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const defaultRandSeed = 1234567890

var (
	defaultRandomVar = rand.New(rand.NewSource(int64(defaultRandSeed)))
)

func newMockPlacementInstance(ctrl *gomock.Controller) services.PlacementInstance {
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

func newMockPlacementInstances(ctrl *gomock.Controller, numInstances int) []services.PlacementInstance {
	svcs := make([]services.PlacementInstance, 0, numInstances)
	for i := 0; i < numInstances; i++ {
		svcs = append(svcs, newMockPlacementInstance(ctrl))
	}
	return svcs
}

func newTestOptions(c *m3em.MockOperatorClient) Options {
	return NewOptions(nil).
		SetOperatorClientFn(func() (*grpc.ClientConn, m3em.OperatorClient, error) {
			return nil, c, nil
		})
}

func TestNodePropertyInitialization(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	opts := newTestOptions(nil)
	mockInstance := newMockPlacementInstance(ctrl)
	serviceNode, err := New(mockInstance, opts)
	require.NoError(t, err)
	require.Equal(t, mockInstance.ID(), serviceNode.ID())
	require.Equal(t, mockInstance.Endpoint(), serviceNode.Endpoint())
	require.Equal(t, mockInstance.Rack(), serviceNode.Rack())
	require.Equal(t, mockInstance.Zone(), serviceNode.Zone())
	require.Equal(t, mockInstance.Weight(), serviceNode.Weight())
	require.Equal(t, mockInstance.Shards(), serviceNode.Shards())
}

func TestNodeErrorStatusIllegalTransitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestOptions(mockClient)
	mockInstance := newMockPlacementInstance(ctrl)
	node, err := New(mockInstance, opts)
	require.NoError(t, err)
	serviceNode := node.(*svcNode)
	require.Equal(t, StatusUninitialized, serviceNode.Status())
	serviceNode.status = StatusError
	require.Error(t, serviceNode.Start())
	require.Error(t, serviceNode.Stop())
	require.Error(t, serviceNode.Setup(nil, nil, "", false))
}

func TestNodeErrorStatusToTeardownTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestOptions(mockClient)
	mockInstance := newMockPlacementInstance(ctrl)
	node, err := New(mockInstance, opts)
	require.NoError(t, err)
	serviceNode := node.(*svcNode)
	require.Equal(t, StatusUninitialized, serviceNode.Status())
	serviceNode.status = StatusError
	mockClient.EXPECT().Teardown(gomock.Any(), gomock.Any())
	require.NoError(t, serviceNode.Teardown())
	require.Equal(t, StatusUninitialized, serviceNode.Status())
}

func TestNodeUninitializedStatusIllegalTransitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestOptions(mockClient)
	mockInstance := newMockPlacementInstance(ctrl)
	node, err := New(mockInstance, opts)
	require.NoError(t, err)
	serviceNode := node.(*svcNode)
	require.Equal(t, StatusUninitialized, serviceNode.Status())
	require.Error(t, serviceNode.Start())
	require.Error(t, serviceNode.Stop())
	require.Error(t, serviceNode.Teardown())
}

func TestNodeUninitializedStatusToSetupTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestOptions(mockClient)
	mb := build.NewMockServiceBuild(ctrl)
	mc := build.NewMockServiceConfiguration(ctrl)
	mockInstance := newMockPlacementInstance(ctrl)
	node, err := New(mockInstance, opts)
	require.NoError(t, err)
	serviceNode := node.(*svcNode)
	require.Equal(t, StatusUninitialized, serviceNode.Status())

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
			Type:       m3em.FileType_SERVICE_BINARY,
			Filename:   "build-id",
			Overwrite:  forceSetup,
			ChunkBytes: dummyBytes,
			ChunkIdx:   0,
		}).Return(nil),
		buildTransferClient.EXPECT().Send(&m3em.TransferRequest{
			Type:       m3em.FileType_SERVICE_BINARY,
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
			Type:       m3em.FileType_SERVICE_CONFIG,
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

	require.NoError(t, serviceNode.Setup(mb, mc, "", forceSetup))
	require.Equal(t, StatusSetup, serviceNode.Status())
	require.Equal(t, mb, serviceNode.currentBuild)
	require.Equal(t, mc, serviceNode.currentConf)
}

func TestNodeSetupStatusIllegalTransitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestOptions(mockClient)
	mockInstance := newMockPlacementInstance(ctrl)
	node, err := New(mockInstance, opts)
	require.NoError(t, err)
	serviceNode := node.(*svcNode)
	serviceNode.status = StatusSetup
	require.Error(t, serviceNode.Stop())
}

func TestNodeSetupStatusToStartTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestOptions(mockClient)
	mockInstance := newMockPlacementInstance(ctrl)
	node, err := New(mockInstance, opts)
	require.NoError(t, err)
	serviceNode := node.(*svcNode)
	serviceNode.status = StatusSetup
	mockClient.EXPECT().Start(gomock.Any(), gomock.Any())
	require.NoError(t, serviceNode.Start())
	require.Equal(t, StatusRunning, serviceNode.Status())
}

func TestNodeSetupStatusToTeardownTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestOptions(mockClient)
	mockInstance := newMockPlacementInstance(ctrl)
	node, err := New(mockInstance, opts)
	require.NoError(t, err)
	serviceNode := node.(*svcNode)
	serviceNode.status = StatusSetup
	mockClient.EXPECT().Teardown(gomock.Any(), gomock.Any())
	require.NoError(t, serviceNode.Teardown())
	require.Equal(t, StatusUninitialized, serviceNode.Status())
}

func TestNodeRunningStatusIllegalTransitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestOptions(mockClient)
	mockInstance := newMockPlacementInstance(ctrl)
	node, err := New(mockInstance, opts)
	require.NoError(t, err)
	serviceNode := node.(*svcNode)
	serviceNode.status = StatusRunning
	require.Error(t, serviceNode.Start())
	require.Error(t, serviceNode.Setup(nil, nil, "", false))
}

func TestNodeRunningStatusToStopTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestOptions(mockClient)
	mockInstance := newMockPlacementInstance(ctrl)
	node, err := New(mockInstance, opts)
	require.NoError(t, err)
	serviceNode := node.(*svcNode)
	serviceNode.status = StatusRunning
	mockClient.EXPECT().Stop(gomock.Any(), gomock.Any())
	require.NoError(t, serviceNode.Stop())
	require.Equal(t, StatusSetup, serviceNode.Status())
}

func TestNodeRunningStatusToTeardownTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := m3em.NewMockOperatorClient(ctrl)
	opts := newTestOptions(mockClient)
	mockInstance := newMockPlacementInstance(ctrl)
	node, err := New(mockInstance, opts)
	require.NoError(t, err)
	serviceNode := node.(*svcNode)
	serviceNode.status = StatusRunning
	mockClient.EXPECT().Teardown(gomock.Any(), gomock.Any())
	require.NoError(t, serviceNode.Teardown())
	require.Equal(t, StatusUninitialized, serviceNode.Status())

}
