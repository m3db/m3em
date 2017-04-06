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
	"fmt"
	"time"

	"github.com/m3db/m3x/instrument"
	xretry "github.com/m3db/m3x/retry"
)

var (
	defaultTimeout            = 2 * time.Minute
	defaultTransferBufferSize = 1024 * 1024 /* 1 MB */
	fourMegaBytes             = 4 * 1024 * 1024
)

type options struct {
	iopts              instrument.Options
	timeout            time.Duration
	transferBufferSize int
	retrier            xretry.Retrier
	hOpts              HeartbeatOptions
}

// NewOptions creates default new options.
func NewOptions(
	opts instrument.Options,
) Options {
	return &options{
		iopts:              opts,
		timeout:            defaultTimeout,
		transferBufferSize: defaultTransferBufferSize,
		retrier:            xretry.NewRetrier(xretry.NewOptions()),
		hOpts:              NewHeartbeatOptions(),
	}
}

func (o *options) Validate() error {
	// grpc max message size (by default, 4MB). see also: https://github.com/grpc/grpc-go/issues/746
	if o.transferBufferSize >= fourMegaBytes {
		return fmt.Errorf("TransferBufferSize must be < 4MB")
	}

	return nil
}

func (o *options) SetInstrumentOptions(io instrument.Options) Options {
	o.iopts = io
	return o
}

func (o *options) InstrumentOptions() instrument.Options {
	return o.iopts
}

func (o *options) SetRetrier(retrier xretry.Retrier) Options {
	o.retrier = retrier
	return o
}

func (o *options) Retrier() xretry.Retrier {
	return o.retrier
}

func (o *options) SetTimeout(d time.Duration) Options {
	o.timeout = d
	return o
}

func (o *options) Timeout() time.Duration {
	return o.timeout
}

func (o *options) SetTransferBufferSize(sz int) Options {
	o.transferBufferSize = sz
	return o
}

func (o *options) TransferBufferSize() int {
	return o.transferBufferSize
}

func (o *options) SetHeartbeatOptions(ho HeartbeatOptions) Options {
	o.hOpts = ho
	return o
}

func (o *options) HeartbeatOptions() HeartbeatOptions {
	return o.hOpts
}
