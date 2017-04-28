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
	hb "github.com/m3db/m3em/generated/proto/heartbeat"
	"github.com/m3db/m3em/generated/proto/m3em"
	mtime "github.com/m3db/m3em/time"

	"github.com/m3db/m3cluster/services"
	"github.com/m3db/m3x/instrument"
	xretry "github.com/m3db/m3x/retry"
	"google.golang.org/grpc"
)

// NodeStatus indicates the different states a ServiceNode can be in. The
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
type NodeStatus int

const (
	// NodeStatusUninitialized refers to the state of an un-initialized node.
	NodeStatusUninitialized NodeStatus = iota

	// NodeStatusSetup is the state of a node which has been Setup()
	NodeStatusSetup

	// NodeStatusRunning is the state of a node which has been Start()-ed
	NodeStatusRunning

	// NodeStatusError is the state of a node which is in an Error state
	NodeStatusError
)

// ServiceNode represents an executable service node. This object controls both the service
// and resources on the host running the service (e.g. fs, processes, etc.)
type ServiceNode interface {
	services.PlacementInstance

	// Setup initializes the directories, config file, and binary for the process being tested.
	// It does not Start the process on the ServiceNode.
	Setup(
		build build.ServiceBuild,
		config build.ServiceConfiguration,
		token string,
		force bool,
	) error

	// Start starts the service process for this ServiceNode.
	Start() error

	// Stop stops the service process for this ServiceNode.
	Stop() error

	// Status returns the ServiceNode status.
	Status() NodeStatus

	// Teardown releases any remote resources used for testing.
	Teardown() error

	// Close releases any locally held resources
	Close() error

	// RegisterListener registers an event listener
	RegisterListener(Listener) ListenerID

	// DeregisterListener un-registers an event listener
	DeregisterListener(ListenerID)

	// Health returns the health for this ServiceNode
	Health() (ServiceNodeHealth, error)

	// TODO(prateek): add more m3db service endpoints in ServiceNode
	// - query service observable properties (nowFn, detailed_status)
	// - set nowFn offset
	// - logs
	// - metrics

	// TODO(prateek): add operator operations for -
	// CleanDataDirectory() error
	// ListDataDirectory(recursive bool, includeContents bool) ([]DirEntry, error)
	// log directory operations
}

// HeartbeatRouter routes heartbeats based on registered servers
type HeartbeatRouter interface {
	hb.HeartbeaterServer

	// Endpoint returns the router endpoint
	Endpoint() string

	// Register registers the specified server under the given id
	Register(string, hb.HeartbeaterServer) error

	// Deregister un-registers any server registered under the given id
	Deregister(string) error
}

// ListenerID is a unique identifier for a registered listener
type ListenerID int

// Listener provides callbacks invoked upon remote process state transitions
type Listener interface {
	// OnProcessTerminate is invoked when the remote process being run terminates
	OnProcessTerminate(node ServiceNode, desc string)

	// OnHeartbeatTimeout is invoked upon remote heartbeats having timed-out
	OnHeartbeatTimeout(node ServiceNode, lastHeartbeatTs time.Time)

	// OnOverwrite is invoked if remote agent control is overwritten by another
	// coordinator
	OnOverwrite(node ServiceNode, desc string)
}

// ServiceNodeHealth provides ServiceNode Health
type ServiceNodeHealth struct {
	Bootstrapped bool
	Status       string
	OK           bool
}

// ServiceNodes is a collection of ServiceNode(s)
type ServiceNodes []ServiceNode

// Options are the knobs used to tweak Environment interactions
type Options interface {
	// Validate validates the Options
	Validate() error

	// SetInstrumentOptions sets the instrumentation options
	SetInstrumentOptions(instrument.Options) Options

	// InstrumentOptions returns the instrumentation options
	InstrumentOptions() instrument.Options

	// TODO(prateek-ref): migrate to ClusterOptions
	// // SetListener sets the ServiceNodeListener
	// SetListener(ServiceNodeListener) Options
	// // Listener returns the ServiceNodeListener
	// Listener() ServiceNodeListener

	// SetNodeOptions sets the NodeOptions
	SetNodeOptions(NodeOptions) Options

	// NodeOptions returns the NodeOptions
	NodeOptions() NodeOptions
}

// NodeOptions are the various knobs to control Node behavior
type NodeOptions interface {
	// Validate validates the NodeOptions
	Validate() error

	// SetInstrumentOptions sets the instrumentation options
	SetInstrumentOptions(instrument.Options) NodeOptions

	// InstrumentOptions returns the instrumentation options
	InstrumentOptions() instrument.Options

	// SetOperationTimeout returns the timeout for node operations
	SetOperationTimeout(time.Duration) NodeOptions

	// OperationTimeout returns the timeout for node operations
	OperationTimeout() time.Duration

	// SetRetrier sets the retrier for node operations
	SetRetrier(xretry.Retrier) NodeOptions

	// OperationRetrier returns the retrier for node operations
	Retrier() xretry.Retrier

	// SetTransferBufferSize sets the bytes buffer size used during file transfer
	SetTransferBufferSize(int) NodeOptions

	// TransferBufferSize returns the bytes buffer size used during file transfer
	TransferBufferSize() int

	// SetHeartbeatOptions sets the HeartbeatOptions
	SetHeartbeatOptions(HeartbeatOptions) NodeOptions

	// HeartbeatOptions returns the HeartbeatOptions
	HeartbeatOptions() HeartbeatOptions

	// SetOperatorClientFn sets the OperatorClientFn
	SetOperatorClientFn(OperatorClientFn) NodeOptions

	// OperatorClientFn returns the OperatorClientFn
	OperatorClientFn() OperatorClientFn
}

// OperatorClientFn returns a function able to construct connections to remote Operators
// TODO(prateek-ref): make this take an input identifier for the node being operated upon
// TODO(prateek-ref): default connection timeout was 2mins
type OperatorClientFn func() (*grpc.ClientConn, m3em.OperatorClient, error)

// HeartbeatOptions are the knobs to control heartbeating behavior
type HeartbeatOptions interface {
	// Validate validates the HeartbeatOptions
	Validate() error

	// SetEnabled sets whether the Heartbeating is enabled
	SetEnabled(bool) HeartbeatOptions

	// Enabled returns whether the Heartbeating is enabled
	Enabled() bool

	// SetNowFn sets the NowFn
	SetNowFn(mtime.NowFn) HeartbeatOptions

	// NowFn returns the NowFn
	NowFn() mtime.NowFn

	// SetInterval sets the heartbeating interval
	SetInterval(time.Duration) HeartbeatOptions

	// Interval returns the heartbeating interval
	Interval() time.Duration

	// SetCheckInterval sets the frequency with which heartbeating timeouts
	// are checked
	SetCheckInterval(time.Duration) HeartbeatOptions

	// CheckInterval returns the frequency with which heartbeating timeouts
	// are checked
	CheckInterval() time.Duration

	// SetTimeout sets the heartbeat timeout duration, i.e. the window of
	// time after which missing heartbeats are considered errorneous
	SetTimeout(time.Duration) HeartbeatOptions

	// Timeout returns the heartbeat timeout duration, i.e. the window of
	// time after which missing heartbeats are considered errorneous
	Timeout() time.Duration

	// SetHeartbeatRouter sets the heartbeat router to be used
	SetHeartbeatRouter(HeartbeatRouter) HeartbeatOptions

	// HeartbeatRouter returns the heartbeat router in use
	HeartbeatRouter() HeartbeatRouter
}
