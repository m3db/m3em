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

	"github.com/m3db/m3x/instrument"
)

const (
	defaultSessionOverride = false
)

type opts struct {
	iopts           instrument.Options
	sessionOverride bool
	token           string
	nodeOpts        NodeOptions
}

// NewOptions returns a new Options object
func NewOptions(iopts instrument.Options) Options {
	if iopts == nil {
		iopts = instrument.NewOptions()
	}
	return &opts{
		iopts:           iopts,
		sessionOverride: defaultSessionOverride,
		nodeOpts:        NewNodeOptions(iopts),
	}
}

func (o opts) Validate() error {
	if o.token == "" {
		return fmt.Errorf("no session token set")
	}
	return o.nodeOpts.Validate()
}

func (o opts) SetInstrumentOptions(iopts instrument.Options) Options {
	o.iopts = iopts
	return o
}

func (o opts) InstrumentOptions() instrument.Options {
	return o.iopts
}

func (o opts) SetSessionToken(t string) Options {
	o.token = t
	return o
}

func (o opts) SessionToken() string {
	return o.token
}

func (o opts) SetSessionOverride(override bool) Options {
	o.sessionOverride = override
	return o
}

func (o opts) SessionOverride() bool {
	return o.sessionOverride
}

func (o opts) SetNodeOptions(no NodeOptions) Options {
	o.nodeOpts = no
	return o
}

func (o opts) NodeOptions() NodeOptions {
	return o.nodeOpts
}
