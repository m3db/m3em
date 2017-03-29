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

type m3dbEnv struct {
	opts        Options
	instances   M3DBInstances
	instanceMap map[string]M3DBInstance
}

// NewM3DBEnvironment returns a new `M3DBEnvironment` of M3DBInstances.
func NewM3DBEnvironment(
	instances []M3DBInstance,
	opts Options,
) (M3DBEnvironment, error) {
	m3dbInstances := make([]M3DBInstance, 0, len(instances))
	m3dbInstanceMap := make(map[string]M3DBInstance, len(instances))
	for _, m3dbInst := range instances {
		m3dbInstances = append(m3dbInstances, m3dbInst)
		m3dbInstanceMap[m3dbInst.ID()] = m3dbInst
	}

	return &m3dbEnv{
		opts:        opts,
		instances:   m3dbInstances,
		instanceMap: m3dbInstanceMap,
	}, nil
}

func (u *m3dbEnv) Instances() M3DBInstances {
	return u.instances
}

func (u *m3dbEnv) InstancesByID() map[string]M3DBInstance {
	return u.instanceMap
}

func (u *m3dbEnv) Status() map[string]InstanceStatus {
	statusMap := make(map[string]InstanceStatus, len(u.instanceMap))
	for id, inst := range u.instanceMap {
		statusMap[id] = inst.Status()
	}
	return statusMap
}
