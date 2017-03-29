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

	"github.com/m3db/m3em/operator"

	"github.com/m3db/m3cluster/services"
	"github.com/m3db/m3x/instrument"
)

const (
	defaultToken                 = ""
	defaultInstanceOverride      = false
	defaultInstanceOpsConcurreny = 3
	defaultInstanceOpRetries     = 3
	defaultInstanceOptTimeout    = 90 * time.Second
)

type opts struct {
	iopts                  instrument.Options
	operatorOpts           operator.Options
	instancePool           string
	serviceID              services.ServiceID
	newInstanceFn          NewInstanceFn
	instanceOpsConcurrency int
	instanceOpRetries      int
	instanceOpTimeout      time.Duration
	instanceOverride       bool
	token                  string
}

// NewOptions returns a new Options object
func NewOptions(iopts instrument.Options) Options {
	if iopts == nil {
		iopts = instrument.NewOptions()
	}
	return &opts{
		iopts:                  iopts,
		instanceOpRetries:      defaultInstanceOpRetries,
		instanceOpsConcurrency: defaultInstanceOpsConcurreny,
		instanceOpTimeout:      defaultInstanceOptTimeout,
		instanceOverride:       defaultInstanceOverride,
		token:                  defaultToken,
	}
}

func (o opts) SetInstrumentOptions(iopts instrument.Options) Options {
	o.iopts = iopts
	return o
}

func (o opts) InstrumentOptions() instrument.Options {
	return o.iopts
}

func (o opts) SetServiceID(sid services.ServiceID) Options {
	o.serviceID = sid
	return o
}

func (o opts) ServiceID() services.ServiceID {
	return o.serviceID
}

func (o opts) SetInstancePool(ip string) Options {
	o.instancePool = ip
	return o
}

func (o opts) InstancePool() string {
	return o.instancePool
}

func (o opts) SetNewInstanceFn(hfn NewInstanceFn) Options {
	o.newInstanceFn = hfn
	return o
}

func (o opts) NewInstanceFn() NewInstanceFn {
	return o.newInstanceFn
}
func (o opts) SetInstanceOperationsConcurrency(c int) Options {
	o.instanceOpsConcurrency = c
	return o
}

func (o opts) InstanceOperationsConcurrency() int {
	return o.instanceOpsConcurrency
}

func (o opts) SetInstanceOperationTimeout(td time.Duration) Options {
	o.instanceOpTimeout = td
	return o
}

func (o opts) InstanceOperationTimeout() time.Duration {
	return o.instanceOpTimeout
}

func (o opts) SetInstanceOperationRetries(r int) Options {
	o.instanceOpRetries = r
	return o
}

func (o opts) InstanceOperationRetries() int {
	return o.instanceOpRetries
}

func (o opts) SetOperatorOptions(oo operator.Options) Options {
	o.operatorOpts = oo
	return o
}

func (o opts) OperatorOptions() operator.Options {
	return o.operatorOpts
}

func (o opts) SetToken(t string) Options {
	o.token = t
	return o
}

func (o opts) Token() string {
	return o.token
}

func (o opts) SetInstanceOverride(override bool) Options {
	o.instanceOverride = override
	return o
}

func (o opts) InstanceOverride() bool {
	return o.instanceOverride
}
