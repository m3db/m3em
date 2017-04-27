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

	"github.com/m3db/m3em/build"

	"github.com/m3db/m3cluster/services"
	"github.com/m3db/m3x/instrument"
)

var (
	defaultSessionOverride = false
	defaultReplication     = 3
	defaultConcurrency     = 10
	defaultNumShards       = 1024
)

type m3dbClusterOpts struct {
	iopts           instrument.Options
	sessionOverride bool
	token           string
	svcBuild        build.ServiceBuild
	svcConf         build.ServiceConfiguration
	placementSvc    services.PlacementService
	maxInstances    int
	replication     int
	numShards       int
	concurrency     int
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
		iopts:           iopts,
		sessionOverride: defaultSessionOverride,
		placementSvc:    placementSvc,
		replication:     defaultReplication,
		numShards:       defaultNumShards,
		concurrency:     defaultConcurrency,
	}
}

func (o m3dbClusterOpts) Validate() error {
	if o.token == "" {
		return fmt.Errorf("no session token set")
	}

	if o.svcBuild == nil {
		return fmt.Errorf("ServiceBuild is not set")
	}

	if o.svcConf == nil {
		return fmt.Errorf("ServiceConf is not set")
	}

	if o.placementSvc == nil {
		return fmt.Errorf("PlacementService is not set")
	}

	return nil
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

func (o m3dbClusterOpts) SetSessionToken(t string) Options {
	o.token = t
	return o
}

func (o m3dbClusterOpts) SessionToken() string {
	return o.token
}

func (o m3dbClusterOpts) SetSessionOverride(override bool) Options {
	o.sessionOverride = override
	return o
}

func (o m3dbClusterOpts) SessionOverride() bool {
	return o.sessionOverride
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
