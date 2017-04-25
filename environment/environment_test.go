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
	"testing"

	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/m3db/m3cluster/services"
	"github.com/stretchr/testify/require"
)

var (
	testHostPool  = "pool"
	testServiceID services.ServiceID
)

func mockM3DBInstances(
	ctrl *gomock.Controller,
	numInstances int,
) M3DBInstances {
	var m3dbInstances M3DBInstances
	for i := 0; i < numInstances; i++ {
		inst := NewMockM3DBInstance(ctrl)
		inst.EXPECT().ID().AnyTimes().Return(fmt.Sprintf("%d", i))
		m3dbInstances = append(m3dbInstances, inst)
	}
	return m3dbInstances
}

func newTestEnv(
	t *testing.T,
	ctrl *gomock.Controller,
	numInstances int,
) (M3DBInstances, M3DBEnvironment) {
	opts := NewOptions(nil)
	instances := mockM3DBInstances(ctrl, numInstances)
	env, err := NewM3DBEnvironment(instances, opts)
	require.NoError(t, err)
	return instances, env
}

func TestNewEnvInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instances, env := newTestEnv(t, ctrl, 5)
	envInstances := env.Instances()
	require.Equal(t, len(instances), len(envInstances))
	for i := range instances {
		require.Equal(t, instances[i], envInstances[i])
	}
}

func TestNewEnvInstancesByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instances, env := newTestEnv(t, ctrl, 5)
	envInstances := env.InstancesByID()
	require.Equal(t, len(instances), len(envInstances))
	for i := range instances {
		id := fmt.Sprintf("%d", i)
		require.Equal(t, instances[i], envInstances[id])
	}
}
