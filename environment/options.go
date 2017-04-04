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

	"github.com/m3db/m3x/instrument"
	xretry "github.com/m3db/m3x/retry"
)

const (
	defaultToken              = ""
	defaultInstanceOverride   = false
	defaultInstanceOptTimeout = 90 * time.Second
)

type opts struct {
	iopts             instrument.Options
	operatorOpts      operator.Options
	instanceOpRetrier xretry.Retrier
	instanceOpTimeout time.Duration
	instanceOverride  bool
	token             string
}

// NewOptions returns a new Options object
func NewOptions(iopts instrument.Options) Options {
	if iopts == nil {
		iopts = instrument.NewOptions()
	}
	return &opts{
		iopts:             iopts,
		instanceOpRetrier: xretry.NewRetrier(xretry.NewOptions()),
		instanceOpTimeout: defaultInstanceOptTimeout,
		instanceOverride:  defaultInstanceOverride,
		token:             defaultToken,
	}
}

func (o opts) SetInstrumentOptions(iopts instrument.Options) Options {
	o.iopts = iopts
	return o
}

func (o opts) InstrumentOptions() instrument.Options {
	return o.iopts
}

func (o opts) SetInstanceOperationTimeout(td time.Duration) Options {
	o.instanceOpTimeout = td
	return o
}

func (o opts) InstanceOperationTimeout() time.Duration {
	return o.instanceOpTimeout
}

func (o opts) SetInstanceOperationRetrier(retrier xretry.Retrier) Options {
	o.instanceOpRetrier = retrier
	return o
}

func (o opts) InstanceOperationRetrier() xretry.Retrier {
	return o.instanceOpRetrier
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
