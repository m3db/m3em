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

	"github.com/m3db/m3cluster/services"
	m3dbclient "github.com/m3db/m3db/client"
	"github.com/m3db/m3x/instrument"
)

var (
	defaultReplication = 3
	defaultConcurrency = 10
	defaultNumShards   = 1024
)

type m3dbClusterOpts struct {
	iopts        instrument.Options
	adminOpts    m3dbclient.AdminOptions
	svcBuild     build.ServiceBuild
	svcConf      build.ServiceConfiguration
	placementSvc services.PlacementService
	maxInstances int
	replication  int
	numShards    int
	concurrency  int
}

// NewOptions returns a new Options object
func NewOptions(
	placementSvc services.PlacementService,
	iopts instrument.Options,
) Options {
	if iopts == nil {
		iopts = instrument.NewOptions()
	}
	return m3dbClusterOpts{
		iopts:        iopts,
		adminOpts:    m3dbclient.NewAdminOptions(),
		placementSvc: placementSvc,
		replication:  defaultReplication,
		numShards:    defaultNumShards,
		concurrency:  defaultConcurrency,
	}
}

func (o m3dbClusterOpts) SetInstrumentOptions(iopts instrument.Options) Options {
	o.iopts = iopts
	return o
}

func (o m3dbClusterOpts) InstrumentOptions() instrument.Options {
	return o.iopts
}

func (o m3dbClusterOpts) SetServiceBuild(b build.ServiceBuild) Options {
	o.svcBuild = b
	return o
}

func (o m3dbClusterOpts) ServiceBuild() build.ServiceBuild {
	return o.svcBuild
}

func (o m3dbClusterOpts) SetServiceConfig(c build.ServiceConfiguration) Options {
	o.svcConf = c
	return o
}

func (o m3dbClusterOpts) ServiceConfig() build.ServiceConfiguration {
	return o.svcConf
}

func (o m3dbClusterOpts) SetMaxInstances(mi int) Options {
	o.maxInstances = mi
	return o
}

func (o m3dbClusterOpts) MaxInstances() int {
	return o.maxInstances
}

func (o m3dbClusterOpts) SetReplication(r int) Options {
	o.replication = r
	return o
}

func (o m3dbClusterOpts) Replication() int {
	return o.replication
}

func (o m3dbClusterOpts) SetNumShards(ns int) Options {
	o.numShards = ns
	return o
}

func (o m3dbClusterOpts) NumShards() int {
	return o.numShards
}

func (o m3dbClusterOpts) SetAdminClientOptions(opts m3dbclient.AdminOptions) Options {
	o.adminOpts = opts
	return o
}

func (o m3dbClusterOpts) AdminClientOptions() m3dbclient.AdminOptions {
	return o.adminOpts
}

func (o m3dbClusterOpts) SetPlacementService(psvc services.PlacementService) Options {
	o.placementSvc = psvc
	return o
}

func (o m3dbClusterOpts) PlacementService() services.PlacementService {
	return o.placementSvc
}

func (o m3dbClusterOpts) SetInstanceConcurrency(c int) Options {
	o.concurrency = c
	return o
}

func (o m3dbClusterOpts) InstanceConcurrency() int {
	return o.concurrency
}
