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

// +build integration

package integration

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/m3db/m3cluster/services/placement"
	"github.com/m3db/m3em/agent"
	"github.com/m3db/m3em/build"
	"github.com/m3db/m3em/environment"
	hb "github.com/m3db/m3em/generated/proto/heartbeat"
	"github.com/m3db/m3em/generated/proto/m3em"

	"github.com/m3db/m3x/instrument"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestHeartbeatTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	targetLocation := newTempDir(t)
	defer os.RemoveAll(targetLocation)

	iopts := instrument.NewOptions()
	// create agent listener, get free port from OS
	agentListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer agentListener.Close()

	// create heartbeat listener, get free port from OS
	heartbeatListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer heartbeatListener.Close()

	// create agent
	aOpts := agent.NewOptions(iopts).SetWorkingDirectory(targetLocation)
	server := grpc.NewServer(grpc.MaxConcurrentStreams(16384))
	agentService, err := agent.New(aOpts)
	require.NoError(t, err)
	m3em.RegisterOperatorServer(server, agentService)
	go func() {
		server.Serve(agentListener)
	}()
	defer server.GracefulStop()

	// create new heartbeat router
	hbRouter := environment.NewHeartbeatRouter(heartbeatListener.Addr().String())
	hbServer := grpc.NewServer(grpc.MaxConcurrentStreams(16384))
	hb.RegisterHeartbeaterServer(hbServer, hbRouter)
	go func() {
		hbServer.Serve(heartbeatListener)
	}()

	// create operator to communicate with agent
	hbOpts := environment.NewHeartbeatOptions().
		SetEnabled(true).
		SetHeartbeatRouter(hbRouter).
		SetTimeout(2 * time.Second).
		SetCheckInterval(100 * time.Millisecond).
		SetInterval(200 * time.Millisecond)
	nodeOpts := environment.NewNodeOptions(iopts).
		SetHeartbeatOptions(hbOpts).
		SetOperatorClientFn(testOperatorClientFn(agentListener.Addr().String()))
	svc := placement.NewInstance()
	node, err := environment.NewServiceNode(svc, nodeOpts)
	require.NoError(t, err)
	defer node.Close()

	// capture notifications from operator
	var (
		lock            sync.Mutex
		receivedTimeout = false
	)
	eventListener := environment.NewListener(
		func(_ environment.ServiceNode, s string) {
			require.FailNow(t, "received termination: %v", s)
		},
		func(_ environment.ServiceNode, ts time.Time) {
			require.False(t, receivedTimeout)
			lock.Lock()
			receivedTimeout = true
			lock.Unlock()
		},
		func(_ environment.ServiceNode, s string) {
			require.FailNow(t, "received overwrite: %v", s)
		},
	)
	id := node.RegisterListener(eventListener)
	defer node.DeregisterListener(id)

	// create test build
	execScript := newTestScript(t, targetLocation, 0, longRunningTestProgram)
	testBuildID := "target-file.sh"
	testBinary := build.NewM3DBBuild(testBuildID, execScript)

	// create test config
	confContents := []byte("some longer string of text\nthat goes on, on and on\n")
	testConfigID := "target-file.conf"
	testConfig := build.NewM3DBConfig(testConfigID, confContents)

	// get the files transferred over
	err = node.Setup(testBinary, testConfig, "tok", false)
	require.NoError(t, err)

	// execute the build
	require.NoError(t, node.Start())

	// stop the heartbeating router so we don't get any more updates
	hbServer.Stop()
	time.Sleep(2 * hbOpts.Timeout())

	// ensure we were notified of exit at the operator level
	require.True(t, receivedTimeout)
	iopts.Logger().Infof("received heartbeat timeout in operator")

	stderrFile := path.Join(targetLocation, fmt.Sprintf("%s.err", testBuildID))
	stderrContents, err := ioutil.ReadFile(stderrFile)
	require.NoError(t, err)
	require.Equal(t, []byte{}, stderrContents, string(stderrContents))

	stdoutFile := path.Join(targetLocation, fmt.Sprintf("%s.out", testBuildID))
	stdoutContents, err := ioutil.ReadFile(stdoutFile)
	require.NoError(t, err)
	require.Equal(t, []byte("testing random output"), stdoutContents)
}
