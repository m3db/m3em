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

package cluster

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/m3db/m3em/build"
	env "github.com/m3db/m3em/environment"
	mockenv "github.com/m3db/m3em/environment/mocks"

	"github.com/golang/mock/gomock"
	"github.com/m3db/m3cluster/services"
	"github.com/stretchr/testify/require"
)

const defaultRandSeed = 1234567890

var (
	defaultRandomVar        = rand.New(rand.NewSource(int64(defaultRandSeed)))
	defaultTestSessionToken = "someLongString"
)

func newDefaultClusterTestOptions(ctrl *gomock.Controller, psvc services.PlacementService) Options {
	mockBuild := build.NewMockServiceBuild(ctrl)
	mockConf := build.NewMockServiceConfiguration(ctrl)
	return NewOptions(psvc, nil).
		SetNumShards(10).
		SetReplication(10).
		SetServiceBuild(mockBuild).
		SetServiceConfig(mockConf).
		SetSessionToken(defaultTestSessionToken)
}

func newMockServiceNode(ctrl *gomock.Controller) env.ServiceNode {
	r := defaultRandomVar
	node := mockenv.NewMockServiceNode(ctrl)
	node.EXPECT().ID().AnyTimes().Return(fmt.Sprintf("%d", r.Int()))
	node.EXPECT().Rack().AnyTimes().Return(fmt.Sprintf("%d", r.Int()))
	node.EXPECT().Endpoint().AnyTimes().Return(fmt.Sprintf("%v:%v", r.Int(), r.Int()))
	node.EXPECT().Zone().AnyTimes().Return(fmt.Sprintf("%d", r.Int()))
	node.EXPECT().Weight().AnyTimes().Return(uint32(r.Int()))
	node.EXPECT().Shards().AnyTimes().Return(nil)
	return node
}

type expectNodeCallTypes struct {
	expectSetup    bool
	expectTeardown bool
	expectStop     bool
	expectStart    bool
}

func addDefaultStatusExpects(nodes []env.ServiceNode, calls expectNodeCallTypes) []env.ServiceNode {
	for _, node := range nodes {
		mNode := node.(*mockenv.MockServiceNode)
		if calls.expectSetup {
			mNode.EXPECT().Setup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
		}
		if calls.expectTeardown {
			mNode.EXPECT().Teardown().AnyTimes().Return(nil)
		}
		if calls.expectStop {
			mNode.EXPECT().Stop().AnyTimes().Return(nil)
		}
		if calls.expectStart {
			mNode.EXPECT().Start().AnyTimes().Return(nil)
		}
	}
	return nodes
}

func newMockServiceNodes(ctrl *gomock.Controller, numNodes int) []env.ServiceNode {
	nodes := make([]env.ServiceNode, 0, numNodes)
	for i := 0; i < numNodes; i++ {
		nodes = append(nodes, newMockServiceNode(ctrl))
	}
	return nodes
}

func newMockPlacementService(ctrl *gomock.Controller) services.PlacementService {
	return services.NewMockPlacementService(ctrl)
}

func TestClusterErrorStatusTransitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPlacementService := newMockPlacementService(ctrl)
	opts := newDefaultClusterTestOptions(ctrl, mockPlacementService)
	expectCalls := expectNodeCallTypes{
		expectSetup:    true,
		expectTeardown: true,
	}
	nodes := addDefaultStatusExpects(newMockServiceNodes(ctrl, 5), expectCalls)
	clusterIface, err := New(nodes, opts)
	require.NoError(t, err)
	cluster := clusterIface.(*svcCluster)
	require.Equal(t, ClusterStatusUninitialized, cluster.Status())
	cluster.status = ClusterStatusError

	// illegal transitions
	require.Error(t, cluster.Setup())
	_, err = cluster.Initialize(1)
	require.Error(t, err)
	require.Error(t, cluster.Start())
	require.Error(t, cluster.StartInitialized())
	require.Error(t, cluster.Stop())
	_, err = cluster.AddNode()
	require.Error(t, err)
	err = cluster.RemoveNode(nil)
	require.Error(t, err)
	_, err = cluster.ReplaceNode(nil)
	require.Error(t, err)

	// teardown (legal)
	require.NoError(t, cluster.Teardown())
	require.Equal(t, ClusterStatusUninitialized, cluster.Status())
}

func TestClusterUninitializedStatusTransitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		mockPlacementService = newMockPlacementService(ctrl)
		mpsvc                = mockPlacementService.(*services.MockPlacementService)
		opts                 = newDefaultClusterTestOptions(ctrl, mockPlacementService)
		expectCalls          = expectNodeCallTypes{expectSetup: true}
		nodes                = addDefaultStatusExpects(newMockServiceNodes(ctrl, 5), expectCalls)
		clusterIface, err    = New(nodes, opts)
	)
	require.NoError(t, err)
	cluster := clusterIface.(*svcCluster)
	require.Equal(t, ClusterStatusUninitialized, cluster.Status())

	// illegal transitions
	require.Error(t, cluster.Start())
	require.Error(t, cluster.StartInitialized())
	require.Error(t, cluster.Stop())
	_, err = cluster.AddNode()
	require.Error(t, err)
	err = cluster.RemoveNode(nil)
	require.Error(t, err)
	_, err = cluster.ReplaceNode(nil)
	require.Error(t, err)

	// setup (legal)
	mpsvc.EXPECT().Placement().Return(nil, 0, nil)
	mpsvc.EXPECT().Delete().Return(nil)
	require.NoError(t, cluster.Setup())
	require.Equal(t, ClusterStatusSetup, cluster.Status())
}

func TestClusterSetupStatusTransitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		mockPlacementService = newMockPlacementService(ctrl)
		mpsvc                = mockPlacementService.(*services.MockPlacementService)
		opts                 = newDefaultClusterTestOptions(ctrl, mockPlacementService)
		expectCalls          = expectNodeCallTypes{
			expectSetup:    true,
			expectTeardown: true,
		}
		nodes        = addDefaultStatusExpects(newMockServiceNodes(ctrl, 5), expectCalls)
		cluster, err = New(nodes, opts)
	)

	require.NoError(t, err)
	require.Equal(t, ClusterStatusUninitialized, cluster.Status())
	mpsvc.EXPECT().Placement().Return(nil, 0, nil)
	mpsvc.EXPECT().Delete().Return(nil)

	require.NoError(t, cluster.Setup())
	require.Equal(t, ClusterStatusSetup, cluster.Status())

	// illegal transitions
	require.Error(t, cluster.Start())
	require.Error(t, cluster.StartInitialized())
	require.Error(t, cluster.Stop())
	_, err = cluster.AddNode()
	require.Error(t, err)
	err = cluster.RemoveNode(nil)
	require.Error(t, err)
	_, err = cluster.ReplaceNode(nil)
	require.Error(t, err)

	// initialize (legal)
	mpsvc.EXPECT().BuildInitialPlacement(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	insts, err := cluster.Initialize(2)
	require.NoError(t, err)
	require.Equal(t, 2, len(insts))

	// teardown (legal)
	mCluster := cluster.(*svcCluster)
	mCluster.status = ClusterStatusSetup
	require.NoError(t, mCluster.Teardown())
	require.Equal(t, ClusterStatusUninitialized, mCluster.Status())
}

func TestClusterRunningStatusTransitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPlacementService := newMockPlacementService(ctrl).(*services.MockPlacementService)
	opts := newDefaultClusterTestOptions(ctrl, mockPlacementService)
	expectCalls := expectNodeCallTypes{
		expectSetup:    true,
		expectTeardown: true,
		expectStop:     true,
	}

	nodes := addDefaultStatusExpects(newMockServiceNodes(ctrl, 5), expectCalls)
	clusterIface, err := New(nodes, opts)
	require.NoError(t, err)
	cluster := clusterIface.(*svcCluster)
	cluster.status = ClusterStatusRunning

	// illegal transitions
	require.Error(t, cluster.Setup())
	require.Error(t, cluster.Start())
	require.Error(t, cluster.StartInitialized())
	_, err = cluster.Initialize(1)
	require.Error(t, err)

	// stop (legal)
	cluster.status = ClusterStatusRunning
	require.NoError(t, cluster.Stop())
	require.Equal(t, ClusterStatusInitialized, cluster.Status())

	// teardown (legal)
	cluster.status = ClusterStatusRunning
	require.NoError(t, cluster.Teardown())
	require.Equal(t, ClusterStatusUninitialized, cluster.Status())

	// add (legal)
	mockPlacementService.EXPECT().AddInstance(gomock.Any()).Return(nil, nil, nil)
	cluster.status = ClusterStatusRunning
	added, err := cluster.AddNode()
	require.NoError(t, err)
	require.Equal(t, ClusterStatusRunning, cluster.Status())

	// replace (legal)
	mockPlacementService.EXPECT().ReplaceInstance(gomock.Any(), gomock.Any()).Return(nil, nil, nil)
	cluster.status = ClusterStatusRunning
	replaced, err := cluster.ReplaceNode(added)
	require.NoError(t, err)
	require.Equal(t, ClusterStatusRunning, cluster.Status())

	// remove (legal)
	mockPlacementService.EXPECT().RemoveInstance(gomock.Any()).Return(nil, nil)
	cluster.status = ClusterStatusRunning
	err = cluster.RemoveNode(replaced)
	require.NoError(t, err)
	require.Equal(t, ClusterStatusRunning, cluster.Status())
}

func TestClusterInitializedStatusTransitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPlacementService := newMockPlacementService(ctrl).(*services.MockPlacementService)
	opts := newDefaultClusterTestOptions(ctrl, mockPlacementService)
	expectCalls := expectNodeCallTypes{
		expectTeardown: true,
		expectStart:    true,
	}

	nodes := addDefaultStatusExpects(newMockServiceNodes(ctrl, 5), expectCalls)
	clusterIface, err := New(nodes, opts)
	require.NoError(t, err)
	cluster := clusterIface.(*svcCluster)
	cluster.status = ClusterStatusInitialized

	// illegal transitions
	require.Error(t, cluster.Setup())
	require.Error(t, cluster.Stop())
	_, err = cluster.Initialize(1)
	require.Error(t, err)

	// start (legal)
	cluster.status = ClusterStatusInitialized
	require.NoError(t, cluster.Start())
	require.Equal(t, ClusterStatusRunning, cluster.Status())

	// StartInitialized (legal)
	cluster.status = ClusterStatusInitialized
	require.NoError(t, cluster.StartInitialized())
	require.Equal(t, ClusterStatusRunning, cluster.Status())

	// teardown (legal)
	cluster.status = ClusterStatusInitialized
	require.NoError(t, cluster.Teardown())
	require.Equal(t, ClusterStatusUninitialized, cluster.Status())

	// add (legal)
	mockPlacementService.EXPECT().AddInstance(gomock.Any()).Return(nil, nil, nil)
	cluster.status = ClusterStatusInitialized
	added, err := cluster.AddNode()
	require.NoError(t, err)
	require.Equal(t, ClusterStatusInitialized, cluster.Status())

	// replace (legal)
	mockPlacementService.EXPECT().ReplaceInstance(gomock.Any(), gomock.Any()).Return(nil, nil, nil)
	cluster.status = ClusterStatusInitialized
	replaced, err := cluster.ReplaceNode(added)
	require.NoError(t, err)
	require.Equal(t, ClusterStatusInitialized, cluster.Status())

	// remove (legal)
	mockPlacementService.EXPECT().RemoveInstance(gomock.Any()).Return(nil, nil)
	cluster.status = ClusterStatusInitialized
	err = cluster.RemoveNode(replaced)
	require.NoError(t, err)
	require.Equal(t, ClusterStatusInitialized, cluster.Status())
}

func TestClusterSetup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		mockPlacementService = newMockPlacementService(ctrl)
		mpsvc                = mockPlacementService.(*services.MockPlacementService)
		opts                 = newDefaultClusterTestOptions(ctrl, mockPlacementService)
		nodes                = newMockServiceNodes(ctrl, 5)
	)
	for _, node := range nodes {
		mi := node.(*mockenv.MockServiceNode)
		mi.EXPECT().Setup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	}
	cluster, err := New(nodes, opts)
	require.NoError(t, err)
	require.Equal(t, ClusterStatusUninitialized, cluster.Status())

	mpsvc.EXPECT().Placement().Return(nil, 0, nil)
	mpsvc.EXPECT().Delete().Return(nil)
	require.NoError(t, cluster.Setup())
	require.Equal(t, ClusterStatusSetup, cluster.Status())
}

func TestClusterInitialize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		mockPlacementService = newMockPlacementService(ctrl)
		mpsvc                = mockPlacementService.(*services.MockPlacementService)
		opts                 = newDefaultClusterTestOptions(ctrl, mockPlacementService)
		nodes                = newMockServiceNodes(ctrl, 5)
	)
	for _, node := range nodes {
		mi := node.(*mockenv.MockServiceNode)
		mi.EXPECT().Setup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	}
	cluster, err := New(nodes, opts)
	require.NoError(t, err)
	require.Equal(t, ClusterStatusUninitialized, cluster.Status())

	mpsvc.EXPECT().Placement().Return(nil, 0, nil)
	mpsvc.EXPECT().Delete().Return(nil)
	require.NoError(t, cluster.Setup())
	require.Equal(t, ClusterStatusSetup, cluster.Status())

	_, err = cluster.Initialize(10)
	require.Error(t, err)

	mpsvc.EXPECT().BuildInitialPlacement(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	usedNodes, err := cluster.Initialize(4)
	require.NoError(t, err)
	require.Equal(t, 4, len(usedNodes))

	spares := cluster.Spares()
	require.Equal(t, 1, len(spares))
	spare := spares[0]
	for _, node := range usedNodes {
		require.NotEqual(t, node.ID(), spare.ID())
	}

	// test StartInitialized
	for _, node := range nodes {
		if spare.ID() != node.ID() {
			mi := node.(*mockenv.MockServiceNode)
			mi.EXPECT().Start().Return(nil)
		}
	}

	require.NoError(t, cluster.StartInitialized())
	require.Equal(t, ClusterStatusRunning, cluster.Status())
}
