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

package exec

import (
	"fmt"
	"os"
	oexec "os/exec"
	"path"
	"path/filepath"
	"sync"
	"time"

	xerrors "github.com/m3db/m3x/errors"
)

var (
	defaultStdoutSuffix         = "out"
	defaultStderrSuffix         = "err"
	defaultProcessCheckInterval = time.Second
)

var (
	errUnableToStartClosed        = fmt.Errorf("unable to start: process monitor Closed()")
	errUnableToStartRunning       = fmt.Errorf("unable to start: process already running")
	errUnableToStopClosed         = fmt.Errorf("unable to stop: process monitor Closed()")
	errUnableToStopStoped         = fmt.Errorf("unable to stop: process not running")
	errUnableToClose              = fmt.Errorf("unable to close: process monitor already Closed()")
	errPathNotSet                 = fmt.Errorf("cmd: no path specified")
	errArgsRequiresPathAsFirstArg = fmt.Errorf("cmd: args[0] un-equal to path")
)

func (c Cmd) validate() error {
	if c.Path == "" {
		return errPathNotSet
	}
	if len(c.Args) > 0 && c.Path != c.Args[0] {
		return errArgsRequiresPathAsFirstArg
	}
	return nil
}

type processListener struct {
	complete func()
	err      func(error)
}

// NewProcessListener returns a new ProcessListener
func NewProcessListener(
	complete func(),
	err func(error),
) ProcessListener {
	return &processListener{
		complete: complete,
		err:      err,
	}
}

func (pl *processListener) OnComplete() {
	if fn := pl.complete; fn != nil {
		fn()
	}
}

func (pl *processListener) OnError(err error) {
	if fn := pl.err; fn != nil {
		fn(err)
	}
}

type startFn func()

type processMonitor struct {
	sync.Mutex
	cmd      *oexec.Cmd
	stdoutFd *os.File
	stderrFd *os.File
	startFn  startFn
	listener ProcessListener
	err      error
	running  bool
	done     bool
}

// NewProcessMonitor creates a new ProcessMonitor
func NewProcessMonitor(cmd Cmd, pl ProcessListener) (ProcessMonitor, error) {
	if err := cmd.validate(); err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(cmd.OutputDir)
	if err != nil {
		return nil, fmt.Errorf("specified directory does not exist: %v", err)
	}

	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("specified path is not a directory: %v", cmd.OutputDir)
	}

	base := filepath.Base(cmd.Path)
	stdoutWriter, err := newWriter(cmd.OutputDir, base, defaultStdoutSuffix)
	if err != nil {
		return nil, fmt.Errorf("unable to open stdout writer: %v", err)
	}

	stderrWriter, err := newWriter(cmd.OutputDir, base, defaultStderrSuffix)
	if err != nil {
		stdoutWriter.Close()
		return nil, fmt.Errorf("unable to open stderr writer: %v", err)
	}

	pm := &processMonitor{
		listener: pl,
		stdoutFd: stdoutWriter,
		stderrFd: stderrWriter,
		cmd: &oexec.Cmd{
			Path:   cmd.Path,
			Args:   cmd.Args,
			Dir:    cmd.OutputDir,
			Stdout: stdoutWriter,
			Stderr: stderrWriter,
		},
	}
	pm.startFn = pm.startAsync

	return pm, nil
}

func (pm *processMonitor) notifyListener(err error) {
	if pm.listener == nil {
		return
	}
	if err == nil {
		pm.listener.OnComplete()
	} else {
		pm.listener.OnError(err)
	}
}

func (pm *processMonitor) Start() error {
	pm.Lock()
	if pm.done {
		pm.Unlock()
		return errUnableToStartClosed
	}

	// TODO(prateek): should we check this?
	if pm.err != nil {
		pm.Unlock()
		return pm.err
	}

	if pm.running {
		pm.Unlock()
		return errUnableToStartRunning
	}
	pm.Unlock()

	pm.startFn()
	return nil
}

func (pm *processMonitor) startAsync() {
	go pm.startAndWait()
}

func (pm *processMonitor) startAndWait() {
	pm.Lock()
	pm.running = true
	pm.Unlock()

	if err := pm.cmd.Start(); err != nil {
		pm.Lock()
		pm.err = err
		pm.running = false
		pm.Unlock()
		pm.notifyListener(err)
		return
	}

	if err := pm.cmd.Wait(); err != nil {
		pm.Lock()
		pm.err = err
		pm.running = false
		pm.Unlock()
		pm.notifyListener(err)
		return
	}

	pm.running = false
	pm.notifyListener(nil)
}

func (pm *processMonitor) Stop() error {
	pm.Lock()
	defer pm.Unlock()
	return pm.stopWithLock()
}

func (pm *processMonitor) stopWithLock() error {
	if pm.done {
		return errUnableToStopClosed
	}

	if pm.err != nil {
		return pm.err
	}

	if !pm.running {
		return errUnableToStopStoped
	}

	pm.running = false
	pm.err = pm.cmd.Process.Kill()
	pm.notifyListener(pm.err)
	return pm.err
}

func (pm *processMonitor) Running() bool {
	pm.Lock()
	defer pm.Unlock()
	return pm.running
}

func (pm *processMonitor) Err() error {
	pm.Lock()
	defer pm.Unlock()
	return pm.err
}

func (pm *processMonitor) Close() error {
	pm.Lock()
	defer pm.Unlock()
	if pm.done {
		return errUnableToClose
	}

	var multiErr xerrors.MultiError
	if pm.running {
		multiErr = multiErr.Add(pm.stopWithLock())
	}

	if pm.stderrFd != nil {
		multiErr = multiErr.Add(pm.stderrFd.Close())
		pm.stderrFd = nil
	}
	if pm.stdoutFd != nil {
		multiErr = multiErr.Add(pm.stdoutFd.Close())
		pm.stdoutFd = nil
	}
	pm.done = true
	return multiErr.FinalError()
}

func newWriter(outputDir string, prefix string, suffix string) (*os.File, error) {
	outputPath := path.Join(outputDir, fmt.Sprintf("%s.%s", prefix, suffix))
	fd, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		return nil, err
	}
	return fd, nil
}
