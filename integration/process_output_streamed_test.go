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
	"testing"
	"time"

	"github.com/m3db/m3em/agent"
	"github.com/m3db/m3em/build"
	hb "github.com/m3db/m3em/generated/proto/heartbeat"
	"github.com/m3db/m3em/generated/proto/m3em"
	"github.com/m3db/m3em/operator"

	"github.com/m3db/m3x/instrument"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestProccessStreamingOutputExecution(t *testing.T) {
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

	// create operator agent
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
	hbRouter := operator.NewHeartbeatRouter(heartbeatListener.Addr().String())
	hbServer := grpc.NewServer(grpc.MaxConcurrentStreams(16384))
	hb.RegisterHeartbeaterServer(hbServer, hbRouter)
	go func() {
		hbServer.Serve(heartbeatListener)
	}()
	defer hbServer.GracefulStop()

	// create operator to communicate with agent
	oOpts := operator.NewOptions(iopts)
	op, err := operator.New(agentListener.Addr().String(), hbRouter, oOpts)
	require.NoError(t, err)

	// create test build
	scriptContents := []byte(`#!/usr/bin/env bash
echo -ne "testing random output"
sleep 100
echo -ne "should never get here"`)
	execScript := newTestScript(t, targetLocation, 0, scriptContents)
	testBuildID := "target-file.sh"
	testBinary := build.NewM3DBBuild(testBuildID, execScript)

	// create test config
	confContents := []byte("some longer string of text\nthat goes on, on and on\n")
	testConfigID := "target-file.conf"
	testConfig := build.NewM3DBConfig(testConfigID, confContents)

	// get the files transferred over
	err = op.Setup(testBinary, testConfig, "tok", false)
	require.NoError(t, err)

	// execute the build
	require.NoError(t, op.Start())
	// ensure the process is running
	stopped := waitUntilAgentFinished(agentService, time.Second)
	require.False(t, stopped)
	require.True(t, agentService.Running())

	// check the output while the process is running to ensure it's not being
	// buffered unreasonably
	stderrFile := path.Join(targetLocation, fmt.Sprintf("%s.err", testBuildID))
	stderrContents, err := ioutil.ReadFile(stderrFile)
	require.NoError(t, err)
	require.Equal(t, []byte{}, stderrContents, string(stderrContents))
	stdoutFile := path.Join(targetLocation, fmt.Sprintf("%s.out", testBuildID))
	stdoutContents, err := ioutil.ReadFile(stdoutFile)
	require.NoError(t, err)
	require.Equal(t, []byte("testing random output"), stdoutContents, string(stdoutContents))

	// process should still be executing
	require.True(t, agentService.Running())
}
