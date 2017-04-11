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
	"time"

	"github.com/m3db/m3em/build"
	"github.com/m3db/m3em/operator"

	"github.com/m3db/m3cluster/services"
	"github.com/m3db/m3x/instrument"
	xretry "github.com/m3db/m3x/retry"
)

// InstanceStatus indicates the different states a ServiceInstance can be in. The
// state diagram below describes the transitions between the various states:
//
//                           ┌──────────────────┐
//                           │                  │
//           ┌Teardown()─────│      Error       │─Reset()─┐
//           │               │                  │         │
//           │               └──────────────────┘         ├───Reset()
//           │                                            │   Setup()
//           ▼                                            ▼     │
// ┌──────────────────┐                         ┌───────────────┴──┐
// │                  │      Setup()            │                  │
// │  Uninitialized   ├────────────────────────▶│      Setup       │◀─┐
// │                  │◀───────────┐            │                  │  │
// └──────────────────┘  Teardown()└────────────└──────────────────┘  │
//           ▲                                            │           │
//           │                                            │           │
//           │                                            │           │
//           │                                  Start()   │           │
//           │                              ┌─────────────┘           │
//           │                              │                         │
//           │                              │                         │
//           │                              │                         │
//           │                              ▼                         │
//           │                    ┌──────────────────┐                │
//           │Teardown()          │                  │            Stop()
//           └────────────────────│     Running      │────────────Reset()
//                                │                  │
//                                └──────────────────┘
type InstanceStatus int

const (
	// InstanceStatusUninitialized refers to the state of an un-initialized instance.
	InstanceStatusUninitialized InstanceStatus = iota

	// InstanceStatusSetup is the state of an instance which has been Setup()
	InstanceStatusSetup

	// InstanceStatusRunning is the state of an instance which has been Start()-ed
	InstanceStatusRunning

	// InstanceStatusError is the state of an instance which is in an Error state
	InstanceStatusError
)

// M3DBInstance represents a testable instance of M3DB. It controls both the service
// and resources on the host running the service (e.g. fs, processes, etc.), the latter is
// available under the Operator() API.
type M3DBInstance interface {
	services.PlacementInstance

	// Setup initializes the directories, config file, and binary for the process being tested.
	// It does not Start the process on the ServiceInstance.
	Setup(build.ServiceBuild, build.ServiceConfiguration) error

	// OverrideConfiguration overrides the instance's service configuration. This only works if
	// the process for the ServiceInstance is not running.
	OverrideConfiguration(build.ServiceConfiguration) error

	// Reset sets the ServiceInstance back to the state after Setup was called.
	Reset() error

	// Teardown releases any resources used for testing.
	Teardown() error

	// Start starts the service process for this ServiceInstance.
	Start() error

	// Stop stops the service process for this ServiceInstance.
	Stop() error

	// Status returns the ServiceInstance status.
	Status() InstanceStatus

	// Operator returns the `Operator` for this ServiceInstance.
	Operator() operator.Operator

	// Health returns the health for this ServiceInstance
	Health() (M3DBInstanceHealth, error)

	// TODO(prateek):
	// - query service observable properties (nowFn, detailed_status)
	// - set nowFn offset
	// - logs
	// - metrics
}

// M3DBInstanceHealth provides M3DBInstance Health
type M3DBInstanceHealth struct {
	Bootstrapped bool
	Status       string
	OK           bool
}

// M3DBInstances is a collection of M3DBInstance(s)
type M3DBInstances []M3DBInstance

// M3DBEnvironment represents a collection of M3DBInstance objects,
// and other resources required to control them.
type M3DBEnvironment interface {
	// Instances returns a list of all the instances in the environment.
	Instances() M3DBInstances

	// InstancesById returns a map [ID -> Instance] of all the
	// instances in the environment.
	InstancesByID() map[string]M3DBInstance

	// Status returns map from instance ID to Status.
	Status() map[string]InstanceStatus
}

// Options are the knobs used to tweak Environment interactions
type Options interface {
	// SetOperatorOptions sets the operator.Options
	SetOperatorOptions(operator.Options) Options

	// OperatorOptions returns the operator.Options
	OperatorOptions() operator.Options

	// SetInstrumentOptions sets the instrumentation options
	SetInstrumentOptions(instrument.Options) Options

	// InstrumentOptions returns the instrumentation options
	InstrumentOptions() instrument.Options

	// SetInstanceOperationTimeout returns the timeout for instance operations
	SetInstanceOperationTimeout(time.Duration) Options

	// InstanceOperationTimeout returns the timeout for instance operations
	InstanceOperationTimeout() time.Duration

	// SetInstanceOperationRetrier sets the retrier for instance operations
	SetInstanceOperationRetrier(xretry.Retrier) Options

	// InstanceOperationRetrier returns the retrier for instance operations
	InstanceOperationRetrier() xretry.Retrier

	// SetSessionToken sets the token used for interactions with remote m3em agents
	SetSessionToken(string) Options

	// SessionToken returns the token used for interactions with remote m3em agents
	SessionToken() string

	// SetSessionOverride sets a flag indicating if m3em agent operations
	// are permitted to override clashing resources
	SetSessionOverride(bool) Options

	// SessionOverride returns a flag indicating if m3em agent operations
	// are permitted to override clashing resources
	SessionOverride() bool
}
