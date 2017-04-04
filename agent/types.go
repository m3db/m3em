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
	"github.com/m3db/m3em/generated/proto/m3em"
	"github.com/m3db/m3em/os/exec"

	"github.com/m3db/m3x/instrument"
)

// Agent is the remote executor of m3em operations
type Agent interface {
	m3em.OperatorServer

	// Running returns a flag indicating if the test process supervised
	// by the Agent is running.
	// TODO(prateek): currently only used in integration tests, could we get rid of it once we have heartbeating?
	Running() bool
}

// Options represent the knobs for a m3em agent
type Options interface {
	// SetInstrumentOptions sets the instrument options
	SetInstrumentOptions(instrument.Options) Options

	// InstrumentOptions returns the instrument options
	InstrumentOptions() instrument.Options

	// SetWorkingDirectory sets the agent's WorkingDirectory
	SetWorkingDirectory(string) Options

	// WorkingDirectory returns the agent's WorkingDirectory
	WorkingDirectory() string

	// SetExecGenFn sets the ExecGenFn
	SetExecGenFn(fn ExecGenFn) Options

	// ExecGenFn returns the ExecGenFn
	ExecGenFn() ExecGenFn

	// SetInitHostResourcesFn sets the InitHostResourcesFn
	SetInitHostResourcesFn(HostResourcesFn) Options

	// InitHostResourcesFn returns the InitHostResourcesFn
	InitHostResourcesFn() HostResourcesFn

	// SetReleaseHostResourcesFn sets the ReleaseHostResourcesFn
	SetReleaseHostResourcesFn(HostResourcesFn) Options

	// ReleaseHostResourcesFn returns the ReleaseHostResourcesFn
	ReleaseHostResourcesFn() HostResourcesFn

	// SetEnvMap sets the EnvMap used to execute any child processes
	SetEnvMap(exec.EnvMap) Options

	// EnvMap returns the EnvMap used to execute any child processes
	EnvMap() exec.EnvMap

	// TODO(prateek): process monitor opts, metric for process uptime
}

// HostResourcesFn is used by the Agent to capture/release any resources
// required on the host. E.g. we use hosts that are typically running
// m3db for staging, for our integration tests as well. So we use this
// function hook to stop any running instances of m3dbnode on the host.
type HostResourcesFn func() error

// ExecGenFn specifies the command to execute for a given build, and config
// e.g. say the process binary expects the config with a cli flag "-f",
// ExecGenFn("binary", "config") == "binary", ["-f", "config"]
type ExecGenFn func(buildPath string, configPath string) (execPath string, args []string)