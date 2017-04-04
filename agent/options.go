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

package agent

import (
	"github.com/m3db/m3x/instrument"
)

var (
	defaultExecGenFn = func(binary string, config string) (string, []string) {
		return binary, []string{"-f", config}
	}

	defaultNoErrorFn = func() error {
		return nil
	}
)

type opts struct {
	iopts      instrument.Options
	workingDir string
	execGenFn  ExecGenFn
	initFn     HostResourcesFn
	releaseFn  HostResourcesFn
}

// NewOptions constructs new options
func NewOptions(io instrument.Options) Options {
	return &opts{
		iopts:     io,
		execGenFn: defaultExecGenFn,
		initFn:    defaultNoErrorFn,
		releaseFn: defaultNoErrorFn,
	}
}

func (o *opts) SetInstrumentOptions(io instrument.Options) Options {
	o.iopts = io
	return o
}

func (o *opts) InstrumentOptions() instrument.Options {
	return o.iopts
}

func (o *opts) SetWorkingDirectory(wd string) Options {
	o.workingDir = wd
	return o
}

func (o *opts) WorkingDirectory() string {
	return o.workingDir
}

func (o *opts) SetExecGenFn(fn ExecGenFn) Options {
	o.execGenFn = fn
	return o
}

func (o *opts) ExecGenFn() ExecGenFn {
	return o.execGenFn
}

func (o *opts) SetInitHostResourcesFn(fn HostResourcesFn) Options {
	o.initFn = fn
	return o
}

func (o *opts) InitHostResourcesFn() HostResourcesFn {
	return o.initFn
}

func (o *opts) SetReleaseHostResourcesFn(fn HostResourcesFn) Options {
	o.releaseFn = fn
	return o
}

func (o *opts) ReleaseHostResourcesFn() HostResourcesFn {
	return o.releaseFn
}
