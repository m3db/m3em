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

package operator

import (
	"time"

	"github.com/m3db/m3x/instrument"
)

var (
	defaultTimeout            = 2 * time.Minute
	defaultTransferBufferSize = 1024 * 1024 /* 1 MB */
)

type options struct {
	iopts              instrument.Options
	timeout            time.Duration
	transferBufferSize int
}

// NewOptions creates default new options.
func NewOptions(
	opts instrument.Options,
) Options {
	return &options{
		iopts:              opts,
		timeout:            defaultTimeout,
		transferBufferSize: defaultTransferBufferSize,
	}
}

func (o *options) SetInstrumentOptions(io instrument.Options) Options {
	o.iopts = io
	return o
}

func (o *options) InstrumentOptions() instrument.Options {
	return o.iopts
}

func (o *options) SetTimeout(d time.Duration) Options {
	o.timeout = d
	return o
}

func (o *options) Timeout() time.Duration {
	return o.timeout
}

// TODO(prateek): ensure this number is less than the grpc max message size (by default, 4MB)
// see also: https://github.com/grpc/grpc-go/issues/746
func (o *options) SetTransferBufferSize(sz int) Options {
	o.transferBufferSize = sz
	return o
}

func (o *options) TransferBufferSize() int {
	return o.transferBufferSize
}
