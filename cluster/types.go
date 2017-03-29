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
	"github.com/m3db/m3em/build"
	env "github.com/m3db/m3em/environment"

	"github.com/m3db/m3cluster/services"
	m3dbc "github.com/m3db/m3db/client"
	"github.com/m3db/m3x/instrument"
)

// TODO(prateek): add 'transitioning' status here

// Status indicates the different states a Cluster can be in. Refer to the
// state diagram below to understand transitions. Note, all states can transition
// to `ClusterStatusError`, these edges are skipped below.
//
//                              ┌───────────────────────────Teardown()────┐
//                              │                               │         │
//                              ▼                               │         │
//                    ╔══════════════════╗                      │         │
//                    ║                  ║                      │         │
//                    ║  Uninitialized   ║─────┐      ╔══════════════════╗│
//                    ║                  ║     │      ║                  ║│
//                    ╚══════════════════╝     │      ║      Error       ║│
//                              ▲              │      ║                  ║│
//                              │              │      ╚══════════════════╝│
//                              ├────┐      Setup()             │         │
//                              │    │         ▼                │         │
//                              │    ╠══════════════════╗       │         │
//                              │    ║                  ║       │         │
//                   Teardown() │    ║      Setup       ║◀────Reset()     │
//                      ┌───────┘    ║                  ║        │        │
//                      │       ┌───▶╚═══╦══════════════╝        │        │
//                      │       │        │     │                 │        │
//                      │     Reset()────┘     │ Initialize()    │        │
//                      │       │              └────────┐        │        │
//                      │       │                       │        │        │
//                      │       │                       │        │        │
//                      │       │                       ▼        │        │
//            ╔═════════╩═══════╩╗   Start()           ╔═════════╩════════╣
//            ║                  ║   StartInitialized()║                  ║
//         ┌──║     Running      ║◀────────────────────║   Initialized    ║──┐
//         │  ║                  ║                     ║                  ║  │
//         │  ╚══════════════════╩───────Stop()───────▶╚══════════════════╝  │
// AddInstance()        ▲                                        ▲           │
// RemoveInstance()     │                                        │     AddInstance()
// ReplaceInstance()────┘                                        └─────RemoveInstance()
//                                                                     ReplaceInstance()
type Status int

const (
	// ClusterStatusUninitialized refers to the state of an un-initialized cluster.
	ClusterStatusUninitialized Status = iota

	// ClusterStatusSetup refers to the state of a cluster whose instances have been
	// setup. In this state, the instances are not running, and the cluster services
	// do not have a defined placement.
	ClusterStatusSetup

	// ClusterStatusInitialized refers to the state of a cluster whose instances have
	// been setup, and the cluster has an assigned placement.
	ClusterStatusInitialized

	// ClusterStatusRunning refers to the state of a cluster with running instances.
	// There is no restriction on the number of instances running, or assigned within
	// the cluster placement.
	ClusterStatusRunning

	// ClusterStatusError refers to a cluster in error.
	ClusterStatusError
)

// Cluster is a collection of clustered M3DB instances.
type Cluster interface {
	// Setup the instances in the Environment provided during construction.
	Setup() error

	// Initialize initializes service placement for the specified numNodes.
	Initialize(numNodes int) ([]env.M3DBInstance, error)

	// AddInstance adds the specified instance to the service placement. It does
	// NOT alter the state of the ServiceInstance (i.e. does not start/stop it).
	AddInstance() (env.M3DBInstance, error)

	// RemoveInstance removes the specified instance from the service placement. It does
	// NOT alter the state of the ServiceInstance (i.e. does not start/stop it).
	RemoveInstance(env.M3DBInstance) error

	// ReplaceInstance replaces the specified instance with a new instance in the service
	// placement. It does NOT alter the state of the TestInstance (i.e. does not start/stop it).
	ReplaceInstance(oldInstance env.M3DBInstance) (env.M3DBInstance, error)

	// Spares returns the instances available in the environment which are not part of the
	// cluster (i.e. placement).
	Spares() []env.M3DBInstance

	// Reset sets the Cluster back to the state after Setup() was called.
	// i.e. it resets placement, and cleans up resources (directories, processes) on the hosts.
	Reset() error

	// Teardown releases the resources acquired during Setup().
	Teardown() error

	// StartInitialized starts any instances which have been initialized and are not running.
	StartInitialized() error

	// Start starts all instances known in the environment, regardless of initialization.
	Start() error

	// Stop stops any running instances in the environment.
	Stop() error

	// Status returns the cluster status
	Status() Status

	// Client returns a m3db Client object
	Client() m3dbc.Client

	// Client returns a m3db AdminClient object
	AdminClient() m3dbc.AdminClient
}

// Options represents the options to configure a `Cluster`
type Options interface {
	// TODO(prateek): this needs a Validate() error

	// SetInstrumentOptions sets the instrumentation options
	SetInstrumentOptions(instrument.Options) Options

	// InstrumentOptions returns the instrumentation options
	InstrumentOptions() instrument.Options

	// SetServiceBuild sets the service build used to configure
	// the cluster
	SetServiceBuild(build.ServiceBuild) Options

	// ServiceBuild returns the service build used in the cluster
	ServiceBuild() build.ServiceBuild

	// SetServiceConfig sets the default service configuration to
	// be used in setting up the cluster
	SetServiceConfig(build.ServiceConfiguration) Options

	// ServiceConfig returns the default service configuration to
	// used in setting up the cluster
	ServiceConfig() build.ServiceConfiguration

	// SetMaxInstances sets the maximum number of instances to be used
	// in the cluster
	SetMaxInstances(int) Options

	// MaxInstances returns the maximum number of instances used in
	// the cluster
	MaxInstances() int

	// SetReplication sets the replication factor for the cluster
	SetReplication(int) Options

	// Replication returns the replication factor by the cluster
	Replication() int

	// SetNumShards sets the number of shards used in the cluster
	SetNumShards(int) Options

	// NumShards returns the number of shards used in the cluster
	NumShards() int

	// SetAdminClientOptions sets the admin client options used in
	// the creation of m3db.Client objects
	SetAdminClientOptions(m3dbc.AdminOptions) Options

	// AdminClientOptions returns the admin client options used in
	// the creation of m3db.Client objects
	AdminClientOptions() m3dbc.AdminOptions

	// SetPlacementService returns the PlacementService to use for cluster
	// configuration
	SetPlacementService(services.PlacementService) Options

	// PlacementService returns the PlacementService to use for cluster
	// configuration
	PlacementService() services.PlacementService
}
