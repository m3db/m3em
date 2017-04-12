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
	m3dbrpc "github.com/m3db/m3db/generated/thrift/rpc"
	m3dbchannel "github.com/m3db/m3db/network/server/tchannelthrift/node/channel"
	tchannel "github.com/uber/tchannel-go"
	"github.com/uber/tchannel-go/thrift"
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
	opts         Options
	id           string
	rack         string
	zone         string
	weight       uint32
	endpoint     string
	shards       shard.Shards
	status       InstanceStatus
	currentBuild build.ServiceBuild
	currentConf  build.ServiceConfiguration
	operator     operator.Operator
	m3dbClient   m3dbrpc.TChanNode
}

// NewM3DBInstance returns a new M3DBInstance.
func NewM3DBInstance(
	inst services.PlacementInstance,
	op operator.Operator,
	opts Options,
) M3DBInstance {
	return &m3dbInst{
		opts:     opts,
		id:       inst.ID(),
		rack:     inst.Rack(),
		zone:     inst.Zone(),
		weight:   inst.Weight(),
		endpoint: inst.Endpoint(),
		shards:   inst.Shards(),
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

	override := i.opts.SessionOverride()
	token := i.opts.SessionToken()
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
	panic("not implemented") // TODO(prateek): implement operator operations for this path
	// return nil
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

func (i *m3dbInst) thriftClient() (m3dbrpc.TChanNode, error) {
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

func (i *m3dbInst) Health() (M3DBInstanceHealth, error) {
	healthResult := M3DBInstanceHealth{}

	client, err := i.thriftClient()
	if err != nil {
		return healthResult, err
	}

	attemptFn := func() error {
		tctx, _ := thrift.NewContext(i.opts.InstanceOperationTimeout())
		result, err := client.Health(tctx)
		if err != nil {
			return err
		}
		healthResult.Bootstrapped = result.GetBootstrapped()
		healthResult.OK = result.GetOk()
		healthResult.Status = result.GetStatus()
		return nil
	}

	retrier := i.opts.InstanceOperationRetrier()
	err = retrier.Attempt(attemptFn)
	return healthResult, err
}
