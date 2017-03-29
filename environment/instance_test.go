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
	"github.com/m3db/m3em/operator"

	"github.com/golang/mock/gomock"
	"github.com/m3db/m3cluster/services"
	"github.com/stretchr/testify/require"
)

const defaultRandSeed = 1234567890

var (
	defaultRandomVar = rand.New(rand.NewSource(int64(defaultRandSeed)))
)

func newMockServiceInstance(ctrl *gomock.Controller) services.PlacementInstance {
	r := defaultRandomVar
	inst := services.NewMockPlacementInstance(ctrl)
	inst.EXPECT().ID().AnyTimes().Return(fmt.Sprintf("%d", r.Int()))
	return inst
}

func newMockOperator(ctrl *gomock.Controller) *operator.MockOperator {
	return operator.NewMockOperator(ctrl)
}

func newMockServiceInstances(ctrl *gomock.Controller, numInstances int) []services.PlacementInstance {
	svcs := make([]services.PlacementInstance, 0, numInstances)
	for i := 0; i < numInstances; i++ {
		svcs = append(svcs, newMockServiceInstance(ctrl))
	}
	return svcs
}

func newDefaultOptions() Options {
	return NewOptions(nil)
}

func TestServiceInstancePropertyInitialization(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	opts := newDefaultOptions()
	mockInstance := newMockServiceInstance(ctrl)
	mockOperator := newMockOperator(ctrl)
	m3dbInstance := NewM3DBInstance(mockInstance.ID(), mockOperator, opts)
	require.Equal(t, mockInstance.ID(), m3dbInstance.ID())
}

func TestServiceInstanceErrorStatusTransitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	opts := newDefaultOptions()
	mockInstance := newMockServiceInstance(ctrl)
	mockOperator := newMockOperator(ctrl)
	m3dbInstance := NewM3DBInstance(mockInstance.ID(), mockOperator, opts).(*m3dbInst)
	require.Equal(t, InstanceStatusUninitialized, m3dbInstance.Status())
	m3dbInstance.status = InstanceStatusError
	require.Error(t, m3dbInstance.Start())
	require.Error(t, m3dbInstance.Stop())
	require.Error(t, m3dbInstance.Setup(nil, nil))

	m3dbInstance.status = InstanceStatusError
	mockOperator.EXPECT().Reset()
	require.NoError(t, m3dbInstance.Reset())
	require.Equal(t, InstanceStatusSetup, m3dbInstance.Status())

	m3dbInstance.status = InstanceStatusError
	mockOperator.EXPECT().Teardown()
	require.NoError(t, m3dbInstance.Teardown())
	require.Equal(t, InstanceStatusUninitialized, m3dbInstance.Status())
}

func TestServiceInstanceStatusTransitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	opts := newDefaultOptions()
	mb := build.NewMockServiceBuild(ctrl)
	mc := build.NewMockServiceConfiguration(ctrl)
	mockInstance := newMockServiceInstance(ctrl)
	mockOperator := newMockOperator(ctrl)
	m3dbInstance := NewM3DBInstance(mockInstance.ID(), mockOperator, opts)
	require.Equal(t, InstanceStatusUninitialized, m3dbInstance.Status())

	// uninitialized -> setup is the only valid (non-errorneous) transition
	require.Error(t, m3dbInstance.Start())
	require.Error(t, m3dbInstance.Stop())
	require.Error(t, m3dbInstance.Reset())
	require.Error(t, m3dbInstance.Teardown())
	mockOperator.EXPECT().Setup(mb, mc, opts.Token(), opts.InstanceOverride())
	require.NoError(t, m3dbInstance.Setup(mb, mc))
	require.Equal(t, InstanceStatusSetup, m3dbInstance.Status())

	// setup -> (running|uninitialized) are the only valid (non-errorneous) transitions
	require.Error(t, m3dbInstance.Stop())
	mockOperator.EXPECT().Setup(mb, mc, opts.Token(), opts.InstanceOverride())
	require.NoError(t, m3dbInstance.Setup(mb, mc))
	mockOperator.EXPECT().Setup(mb, mc, opts.Token(), opts.InstanceOverride())
	require.NoError(t, m3dbInstance.Setup(mb, mc))
	mockOperator.EXPECT().Reset()
	require.NoError(t, m3dbInstance.Reset())
	require.Equal(t, InstanceStatusSetup, m3dbInstance.Status())
	mockOperator.EXPECT().Teardown()
	require.NoError(t, m3dbInstance.Teardown())
	require.Equal(t, InstanceStatusUninitialized, m3dbInstance.Status())
	mockOperator.EXPECT().Setup(mb, mc, opts.Token(), opts.InstanceOverride())
	require.NoError(t, m3dbInstance.Setup(mb, mc))
	require.Equal(t, InstanceStatusSetup, m3dbInstance.Status())
	mockOperator.EXPECT().Start()
	require.NoError(t, m3dbInstance.Start())
	require.Equal(t, InstanceStatusRunning, m3dbInstance.Status())

	// running -> (setup|uninitialized) are the only valid (non-errorneous) transitions
	require.Error(t, m3dbInstance.Start())
	require.Error(t, m3dbInstance.Setup(mb, mc))
	mockOperator.EXPECT().Teardown()
	require.NoError(t, m3dbInstance.Teardown())
	require.Equal(t, InstanceStatusUninitialized, m3dbInstance.Status())
	mockOperator.EXPECT().Setup(mb, mc, opts.Token(), opts.InstanceOverride())
	require.NoError(t, m3dbInstance.Setup(mb, mc))
	require.Equal(t, InstanceStatusSetup, m3dbInstance.Status())
	mockOperator.EXPECT().Start()
	require.NoError(t, m3dbInstance.Start())
	require.Equal(t, InstanceStatusRunning, m3dbInstance.Status())
	mockOperator.EXPECT().Stop()
	require.NoError(t, m3dbInstance.Stop())
	require.Equal(t, InstanceStatusSetup, m3dbInstance.Status())
	mockOperator.EXPECT().Start()
	require.NoError(t, m3dbInstance.Start())
	require.Equal(t, InstanceStatusRunning, m3dbInstance.Status())
	mockOperator.EXPECT().Reset()
	require.NoError(t, m3dbInstance.Reset())
	require.Equal(t, InstanceStatusSetup, m3dbInstance.Status())
}

func TestServiceInstanceSetupBuildVerify(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	opts := newDefaultOptions()
	mb := build.NewMockServiceBuild(ctrl)
	mc := build.NewMockServiceConfiguration(ctrl)
	mockInstance := newMockServiceInstance(ctrl)
	mockOperator := newMockOperator(ctrl)
	m3dbInstance := NewM3DBInstance(mockInstance.ID(), mockOperator, opts)
	require.Equal(t, InstanceStatusUninitialized, m3dbInstance.Status())

	mockOperator.EXPECT().Setup(mb, mc, opts.Token(), opts.InstanceOverride())
	require.NoError(t, m3dbInstance.Setup(mb, mc))
	require.Equal(t, InstanceStatusSetup, m3dbInstance.Status())

	inst := m3dbInstance.(*m3dbInst)
	require.Equal(t, mb, inst.currentBuild)
	require.Equal(t, mc, inst.currentConf)

	mcp := build.NewMockServiceConfiguration(ctrl)
	require.NoError(t, m3dbInstance.OverrideConfiguration(mcp))
	require.Equal(t, mcp, inst.currentConf)
}
