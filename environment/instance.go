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
	"sync"

	"github.com/m3db/m3em/build"
	"github.com/m3db/m3em/operator"

	"github.com/m3db/m3cluster/services"
	"github.com/m3db/m3cluster/shard"
)

var (
	errUnableToSetupInitializedInstance = fmt.Errorf("unable to setup instance, must be either setup/uninitialized")
	errUnableToResetInstance            = fmt.Errorf("unable to reset instance, must be either setup/running")
	errUnableToTeardownInstance         = fmt.Errorf("unable to teardown instance, must be either setup/running")
	errUnableToStartInstance            = fmt.Errorf("unable to start instance, it must be setup")
	errUnableToStopInstance             = fmt.Errorf("unable to stop instance, it must be running")
	errUnableToOverrideConf             = fmt.Errorf("unable to override instance configuration, it must be setup")
)

type m3dbInst struct {
	sync.Mutex
	id           string
	rack         string
	zone         string
	weight       uint32
	endpoint     string
	shards       shard.Shards
	operator     operator.Operator
	opts         Options
	currentBuild build.ServiceBuild
	currentConf  build.ServiceConfiguration
	status       InstanceStatus
}

// NewM3DBInstance returns a new M3DBInstance.
func NewM3DBInstance(
	id string,
	op operator.Operator,
	opts Options,
) M3DBInstance {
	return &m3dbInst{
		opts:     opts,
		id:       id,
		status:   InstanceStatusUninitialized,
		operator: op,
	}
}

func (i *m3dbInst) String() string {
	i.Lock()
	defer i.Unlock()
	return fmt.Sprintf(
		"Instance[ID=%s, Rack=%s, Zone=%s, Weight=%d, Endpoint=%s, Shards=%s]",
		i.id, i.rack, i.zone, i.weight, i.endpoint, i.shards.String(),
	)
}

func (i *m3dbInst) ID() string {
	i.Lock()
	defer i.Unlock()
	return i.id
}

func (i *m3dbInst) SetID(id string) services.PlacementInstance {
	i.Lock()
	defer i.Unlock()
	i.id = id
	return i
}

func (i *m3dbInst) Rack() string {
	i.Lock()
	defer i.Unlock()
	return i.rack
}

func (i *m3dbInst) SetRack(r string) services.PlacementInstance {
	i.Lock()
	defer i.Unlock()
	i.rack = r
	return i
}

func (i *m3dbInst) Zone() string {
	i.Lock()
	defer i.Unlock()
	return i.zone
}

func (i *m3dbInst) SetZone(z string) services.PlacementInstance {
	i.Lock()
	defer i.Unlock()
	i.zone = z
	return i
}

func (i *m3dbInst) Weight() uint32 {
	i.Lock()
	defer i.Unlock()
	return i.weight
}

func (i *m3dbInst) SetWeight(w uint32) services.PlacementInstance {
	i.Lock()
	defer i.Unlock()
	i.weight = w
	return i
}

func (i *m3dbInst) Endpoint() string {
	i.Lock()
	defer i.Unlock()
	return i.endpoint
}

func (i *m3dbInst) SetEndpoint(ip string) services.PlacementInstance {
	i.Lock()
	defer i.Unlock()
	i.endpoint = ip
	return i
}

func (i *m3dbInst) Shards() shard.Shards {
	i.Lock()
	defer i.Unlock()
	return i.shards
}

func (i *m3dbInst) SetShards(s shard.Shards) services.PlacementInstance {
	i.Lock()
	defer i.Unlock()
	i.shards = s
	return i
}

func (i *m3dbInst) Setup(b build.ServiceBuild, c build.ServiceConfiguration) error {
	i.Lock()
	defer i.Unlock()
	if i.status != InstanceStatusUninitialized &&
		i.status != InstanceStatusSetup {
		return errUnableToSetupInitializedInstance
	}

	override := i.opts.InstanceOverride()
	token := i.opts.Token()
	i.currentConf = c
	i.currentBuild = b

	if err := i.operator.Setup(b, c, token, override); err != nil {
		return err
	}

	i.status = InstanceStatusSetup
	return nil
}

func (i *m3dbInst) OverrideConfiguration(conf build.ServiceConfiguration) error {
	i.Lock()
	defer i.Unlock()
	if i.status != InstanceStatusSetup {
		return errUnableToOverrideConf
	}
	i.currentConf = conf
	// TODO(prateek): implement operator operations for this path
	return nil
}

func (i *m3dbInst) Reset() error {
	i.Lock()
	defer i.Unlock()
	if status := i.status; status != InstanceStatusRunning &&
		status != InstanceStatusSetup && status != InstanceStatusError {
		return errUnableToResetInstance
	}

	if err := i.operator.Reset(); err != nil {
		return err
	}

	i.status = InstanceStatusSetup
	return nil
}

func (i *m3dbInst) Teardown() error {
	i.Lock()
	defer i.Unlock()
	if status := i.status; status != InstanceStatusRunning &&
		status != InstanceStatusSetup &&
		status != InstanceStatusError {
		return errUnableToTeardownInstance
	}

	if err := i.operator.Teardown(); err != nil {
		return err
	}

	i.status = InstanceStatusUninitialized
	return nil
}

func (i *m3dbInst) Start() error {
	i.Lock()
	defer i.Unlock()
	if i.status != InstanceStatusSetup {
		return errUnableToStartInstance
	}

	if err := i.operator.Start(); err != nil {
		return err
	}

	i.status = InstanceStatusRunning
	return nil
}

func (i *m3dbInst) Stop() error {
	i.Lock()
	defer i.Unlock()
	if i.status != InstanceStatusRunning {
		return errUnableToStopInstance
	}

	if err := i.operator.Stop(); err != nil {
		return err
	}

	i.status = InstanceStatusSetup
	return nil
}

func (i *m3dbInst) Status() InstanceStatus {
	i.Lock()
	defer i.Unlock()
	return i.status
}

func (i *m3dbInst) Operator() operator.Operator {
	return i.operator
}
