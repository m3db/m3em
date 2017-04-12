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
	"context"
	"sync"
	"testing"
	"time"

	hb "github.com/m3db/m3em/generated/proto/heartbeat"

	"github.com/golang/mock/gomock"
	"github.com/m3db/m3x/instrument"
	"github.com/stretchr/testify/require"
)

func newTestHeartbeatOpts() HeartbeatOptions {
	return NewHeartbeatOptions().
		SetCheckInterval(10 * time.Millisecond).
		SetInterval(20 * time.Millisecond).
		SetTimeout(100 * time.Millisecond)
}

func newTestListener(t *testing.T) *listener {
	return &listener{
		onProcessTerminate: func(desc string) { require.Fail(t, "onProcessTerminate invoked %s", desc) },
		onHeartbeatTimeout: func(ts time.Time) { require.Fail(t, "onHeartbeatTimeout invoked %s", ts.String()) },
		onOverwrite:        func(desc string) { require.Fail(t, "onOverwrite invoked %s", desc) },
	}
}

func TestHeartbeaterSimple(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		lg       = newListenerGroup()
		opts     = newTestHeartbeatOpts()
		iopts    = instrument.NewOptions()
		hbServer = newHeartbeater(lg, opts, iopts)
	)
	require.NoError(t, hbServer.start())

	hbServer.Heartbeat(context.Background(),
		&hb.HeartbeatRequest{
			Code:           hb.HeartbeatCode_HEALTHY,
			ProcessRunning: false,
		})
	time.Sleep(10 * time.Millisecond) // to yield to any pending go routines
	hbServer.stop()
}

func TestHeartbeatingUnknownCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		lg       = newListenerGroup()
		opts     = newTestHeartbeatOpts()
		iopts    = instrument.NewOptions()
		hbServer = newHeartbeater(lg, opts, iopts)
	)
	require.NoError(t, hbServer.start())

	_, err := hbServer.Heartbeat(context.Background(),
		&hb.HeartbeatRequest{
			Code:           hb.HeartbeatCode_UNKNOWN,
			ProcessRunning: false,
		})
	require.Error(t, err)
	time.Sleep(10 * time.Millisecond) // to yield to any pending go routines
	hbServer.stop()
}

func TestHeartbeatingProcessTermination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		lg    = newListenerGroup()
		opts  = newTestHeartbeatOpts()
		iopts = instrument.NewOptions()
	)

	var (
		lock              sync.Mutex
		processTerminated = false
		lnr               = newTestListener(t)
	)
	lnr.onProcessTerminate = func(string) {
		lock.Lock()
		processTerminated = true
		lock.Unlock()
	}
	lg.add(lnr)

	hbServer := newHeartbeater(lg, opts, iopts)
	require.NoError(t, hbServer.start())
	_, err := hbServer.Heartbeat(context.Background(),
		&hb.HeartbeatRequest{
			Code: hb.HeartbeatCode_PROCESS_TERMINATION,
		})
	require.NoError(t, err)
	time.Sleep(10 * time.Millisecond) // to yield to any pending go routines
	hbServer.stop()

	lock.Lock()
	defer lock.Unlock()
	require.True(t, processTerminated)
}

func TestHeartbeatingOverwrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		lg    = newListenerGroup()
		opts  = newTestHeartbeatOpts()
		iopts = instrument.NewOptions()
	)

	var (
		lock        sync.Mutex
		overwritten = false
		lnr         = newTestListener(t)
	)
	lnr.onOverwrite = func(string) {
		lock.Lock()
		overwritten = true
		lock.Unlock()
	}
	lg.add(lnr)

	hbServer := newHeartbeater(lg, opts, iopts)
	require.NoError(t, hbServer.start())
	_, err := hbServer.Heartbeat(context.Background(),
		&hb.HeartbeatRequest{
			Code: hb.HeartbeatCode_OVERWRITTEN,
		})
	require.NoError(t, err)
	time.Sleep(10 * time.Millisecond) // to yield to any pending go routines
	hbServer.stop()

	lock.Lock()
	defer lock.Unlock()
	require.True(t, overwritten)
}

func TestHeartbeatingTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		now       = time.Now()
		callCount = 0
		nowFn     = func() time.Time {
			callCount++
			if callCount == 1 {
				return now
			}
			return now.Add(time.Hour)
		}
		lg    = newListenerGroup()
		opts  = newTestHeartbeatOpts().SetNowFn(nowFn)
		iopts = instrument.NewOptions()
	)

	var (
		lock     sync.Mutex
		timedout = false
		lnr      = newTestListener(t)
	)
	lnr.onHeartbeatTimeout = func(ts time.Time) {
		lock.Lock()
		timedout = true
		lock.Unlock()
	}
	lg.add(lnr)

	hbServer := newHeartbeater(lg, opts, iopts)
	require.NoError(t, hbServer.start())
	_, err := hbServer.Heartbeat(context.Background(),
		&hb.HeartbeatRequest{
			Code: hb.HeartbeatCode_HEALTHY,
		})
	require.NoError(t, err)
	time.Sleep(10 * time.Millisecond)
	hbServer.stop()

	lock.Lock()
	defer lock.Unlock()
	require.True(t, timedout)
}
