package operator

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/m3db/m3em/generated/proto/m3em"

	"github.com/golang/mock/gomock"
	"github.com/m3db/m3x/instrument"
	"github.com/stretchr/testify/require"
)

func newTestHeartbeatOpts() HeartbeatOptions {
	return NewHeartbeatOptions().
		SetCheckInterval(100 * time.Millisecond).
		SetInterval(2 * time.Second).
		SetTimeout(10 * time.Second)
}

func newTestListener(t *testing.T) *listener {
	return &listener{
		onProcessTerminate: func(desc string) { require.Fail(t, "onProcessTerminate invoked %s", desc) },
		onHeartbeatTimeout: func(ts time.Time) { require.Fail(t, "onHeartbeatTimeout invoked %s", ts.String()) },
		onOverwrite:        func(desc string) { require.Fail(t, "onOverwrite invoked %s", desc) },
		onInternalError:    func(err error) { require.Fail(t, "onInternalError invoked %s", err.Error()) },
	}
}

func TestHeartbeaterCancellation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		lg      = newListenerGroup()
		opts    = newTestHeartbeatOpts()
		iopts   = instrument.NewOptions()
		client  = m3em.NewMockOperatorClient(ctrl)
		hClient = m3em.NewMockOperator_HeartbeatClient(ctrl)
	)
	hClient.EXPECT().Recv().Return(&m3em.HeartbeatResponse{
		Code:           m3em.HeartbeatCode_HEARTBEAT_CODE_HEALTHY,
		ProcessRunning: false,
	}, nil).AnyTimes()

	client.EXPECT().Heartbeat(gomock.Any(), &m3em.HeartbeatRequest{
		FrequencySecs: uint32(opts.Interval().Seconds()),
	}).Return(hClient, nil)

	hb := newHeartbeater(client, lg, iopts, opts)
	require.NoError(t, hb.start())

	notCalled := true
	hb.clientCancel = func() {
		hClient = m3em.NewMockOperator_HeartbeatClient(ctrl)
		hClient.EXPECT().Recv().Return(nil, context.Canceled)
		hb.client = hClient
		notCalled = false
	}

	time.Sleep(10 * time.Millisecond) // to yield to any pending go routines
	hb.stop()
	require.False(t, notCalled)
}

func TestHeartbeatingUnknownCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		lg      = newListenerGroup()
		opts    = newTestHeartbeatOpts()
		iopts   = instrument.NewOptions()
		client  = m3em.NewMockOperatorClient(ctrl)
		hClient = m3em.NewMockOperator_HeartbeatClient(ctrl)
	)

	var (
		lock                  sync.Mutex
		internalErrorNotified = false
		lnr                   = newTestListener(t)
	)
	lnr.onInternalError = func(err error) {
		println(err.Error())
		lock.Lock()
		internalErrorNotified = true
		lock.Unlock()
	}
	lg.add(lnr)

	hClient.EXPECT().Recv().Return(&m3em.HeartbeatResponse{
		Code: m3em.HeartbeatCode_HEARTBEAT_CODE_UNKNOWN,
	}, nil)

	client.EXPECT().Heartbeat(gomock.Any(), &m3em.HeartbeatRequest{
		FrequencySecs: uint32(opts.Interval().Seconds()),
	}).Return(hClient, nil)

	hb := newHeartbeater(client, lg, iopts, opts)
	require.NoError(t, hb.start())
	time.Sleep(10 * time.Millisecond) // to yield to any pending go routines
	lock.Lock()
	defer lock.Unlock()
	require.True(t, internalErrorNotified)
}

func TestHeartbeatingEOF(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		lg      = newListenerGroup()
		opts    = newTestHeartbeatOpts()
		iopts   = instrument.NewOptions()
		client  = m3em.NewMockOperatorClient(ctrl)
		hClient = m3em.NewMockOperator_HeartbeatClient(ctrl)
	)

	var (
		lock                  sync.Mutex
		internalErrorNotified = false
		lnr                   = newTestListener(t)
	)
	lnr.onInternalError = func(err error) {
		lock.Lock()
		internalErrorNotified = true
		lock.Unlock()
	}
	lg.add(lnr)

	hClient.EXPECT().Recv().Return(nil, io.EOF)
	client.EXPECT().Heartbeat(gomock.Any(), &m3em.HeartbeatRequest{
		FrequencySecs: uint32(opts.Interval().Seconds()),
	}).Return(hClient, nil)

	hb := newHeartbeater(client, lg, iopts, opts)
	require.NoError(t, hb.start())
	time.Sleep(10 * time.Millisecond) // to yield to any pending go routines
	lock.Lock()
	defer lock.Unlock()
	require.True(t, internalErrorNotified)
}

func TestHeartbeatingProcessTermination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		lg      = newListenerGroup()
		opts    = newTestHeartbeatOpts()
		iopts   = instrument.NewOptions()
		client  = m3em.NewMockOperatorClient(ctrl)
		hClient = m3em.NewMockOperator_HeartbeatClient(ctrl)
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

	hClient.EXPECT().Recv().Do(func() {
		time.Sleep(time.Hour)
	}).After(
		hClient.EXPECT().Recv().Return(&m3em.HeartbeatResponse{
			Code: m3em.HeartbeatCode_HEARTBEAT_CODE_PROCESS_TERMINATION,
		}, nil),
	)
	client.EXPECT().Heartbeat(gomock.Any(), &m3em.HeartbeatRequest{
		FrequencySecs: uint32(opts.Interval().Seconds()),
	}).Return(hClient, nil)

	hb := newHeartbeater(client, lg, iopts, opts)
	require.NoError(t, hb.start())
	time.Sleep(10 * time.Millisecond) // to yield to any pending go routines
	lock.Lock()
	defer lock.Unlock()
	require.True(t, processTerminated)
}

func TestHeartbeatingOverwrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		lg      = newListenerGroup()
		opts    = newTestHeartbeatOpts()
		iopts   = instrument.NewOptions()
		client  = m3em.NewMockOperatorClient(ctrl)
		hClient = m3em.NewMockOperator_HeartbeatClient(ctrl)
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

	hClient.EXPECT().Recv().Return(&m3em.HeartbeatResponse{
		Code: m3em.HeartbeatCode_HEARTBEAT_CODE_OVERWRITTEN,
	}, nil)
	client.EXPECT().Heartbeat(gomock.Any(), &m3em.HeartbeatRequest{
		FrequencySecs: uint32(opts.Interval().Seconds()),
	}).Return(hClient, nil)

	hb := newHeartbeater(client, lg, iopts, opts)
	require.NoError(t, hb.start())
	time.Sleep(10 * time.Millisecond) // to yield to any pending go routines
	lock.Lock()
	defer lock.Unlock()
	require.True(t, overwritten)
}
