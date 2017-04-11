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

	"github.com/m3db/m3em/build"

	"github.com/m3db/m3x/instrument"
	xretry "github.com/m3db/m3x/retry"
)

// Operator provides control over the resources (fs/process/etc) available to a ServiceInstance.
type Operator interface {
	// Setup sets up the remote host with resources required to execute a test build
	// NB(prateek): `force` is used to determine if any existing resources (e.g. existing files, running
	// processes) on the remote host may be trampled.
	Setup(
		build build.ServiceBuild,
		config build.ServiceConfiguration,
		token string,
		force bool,
	) error

	// Teardown releases resources acquired on the remote host
	Teardown() error

	// Start begins the execution of the test build
	Start() error

	// Stop terminates the execution of the test build
	Stop() error

	// Reset gets the remote host to same state as after Setup()
	Reset() error

	// RegisterListener registers an event listener
	RegisterListener(Listener) ListenerID

	// DeregisterListener un-registers an event listener
	DeregisterListener(ListenerID)

	// CleanDataDirectory() error
	// ListDataDirectory(recursive bool, includeContents bool) ([]DirEntry, error)

	// log directory operations
}

// ListenerID is a unique identifier for a registered listener
type ListenerID int

// Listener provides callbacks invoked upon remote process state transitions
type Listener interface {
	// OnProcessTerminate is invoked when the remote process being run terminates
	OnProcessTerminate(desc string)

	// OnHeartbeatTimeout is invoked upon remote heartbeats having timed-out
	OnHeartbeatTimeout(lastHeartbeatTs time.Time)

	// OnOverwrite is invoked if remote agent control is overwritten by another
	// coordinator
	OnOverwrite(desc string)
}

// Options returns options pertaining to `Operator` configuration.
type Options interface {
	// Validate returns an error iff the options are set to invalid values
	Validate() error

	// SetInstrumentOptions sets the instrumentation options
	SetInstrumentOptions(instrument.Options) Options

	// InstrumentOptions returns the instrumentation options
	InstrumentOptions() instrument.Options

	// SetRetrier sets the retrier
	SetRetrier(xretry.Retrier) Options

	// Retrier returns the retrier
	Retrier() xretry.Retrier

	// SetTimeout sets the GRPC connection timeout
	SetTimeout(time.Duration) Options

	// Timeout returns the GRPC connection timeout
	Timeout() time.Duration

	// SetTransferBufferSize sets the bytes buffer size used during file transfer
	SetTransferBufferSize(int) Options

	// TransferBufferSize returns the bytes buffer size used during file transfer
	TransferBufferSize() int

	// SetHeartbeatOptions sets the HeartbeatOptions
	SetHeartbeatOptions(HeartbeatOptions) Options

	// HeartbeatOptions returns the HeartbeatOptions
	HeartbeatOptions() HeartbeatOptions
}

// HeartbeatOptions are the knobs to control heartbeating behavior
type HeartbeatOptions interface {
	// SetEnabled sets whether the Heartbeating is enabled
	SetEnabled(bool) HeartbeatOptions

	// Enabled returns whether the Heartbeating is enabled
	Enabled() bool

	// SetNowFn sets the NowFn
	SetNowFn(NowFn) HeartbeatOptions

	// NowFn returns the NowFn
	NowFn() NowFn

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
}

// NowFn returns the current time
type NowFn func() time.Time

// DirEntry corresponds to a file present in a directory.
type DirEntry struct {
	FilePath     string
	FileContents []byte // TODO(prateek): make this FileReaderIter
}
