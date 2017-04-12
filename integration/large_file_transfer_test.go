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
	"io/ioutil"
	"net"
	"os"
	"path"
	"testing"

	"github.com/m3db/m3em/agent"
	"github.com/m3db/m3em/build"
	"github.com/m3db/m3em/generated/proto/m3em"
	"github.com/m3db/m3em/operator"

	"github.com/m3db/m3x/instrument"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestLargeFileTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	targetLocation := newTempDir(t)
	defer os.RemoveAll(targetLocation)

	iopts := instrument.NewOptions()
	// create listener, get free port from OS
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	defer l.Close()

	// create operator agent
	aOpts := agent.NewOptions(iopts).SetWorkingDirectory(targetLocation)
	server := grpc.NewServer(grpc.MaxConcurrentStreams(16384))
	service, err := agent.New(aOpts)
	require.NoError(t, err)
	m3em.RegisterOperatorServer(server, service)
	go func() {
		server.Serve(l)
	}()
	defer server.GracefulStop()

	// create operator to communicate with agent
	oOpts := operator.NewOptions(iopts)
	op, err := operator.New(l.Addr().String(), oOpts)
	require.NoError(t, err)

	// create test build
	largeTestFileSize := int64(1024 * 1024 * 100) /* 100 MB */
	largeTestFile := newLargeTempFile(t, targetLocation, largeTestFileSize)
	defer os.Remove(largeTestFile.Name())
	testBuildID := "target-file.out"
	targetBuildFile := path.Join(targetLocation, testBuildID)
	testBinary := build.NewM3DBBuild(testBuildID, largeTestFile.Name())

	// create test config
	confContents := []byte("some longer string of text\nthat goes on, on and on\n")
	testConfigID := "target-file.conf"
	targetConfigFile := path.Join(targetLocation, testConfigID)
	testConfig := build.NewM3DBConfig(testConfigID, confContents)

	err = op.Setup(testBinary, testConfig, "tok", false)
	require.NoError(t, err)

	// test copied build file contents
	obsBytes, err := ioutil.ReadFile(targetBuildFile)
	require.NoError(t, err)
	// NB(prateek): not testing if bytes are equal as this is a large file,
	// the operator <-> agent communication verifies checksum and num chunks
	// sent are correct.

	// test copied config file contents
	obsBytes, err = ioutil.ReadFile(targetConfigFile)
	require.NoError(t, err)
	require.Equal(t, confContents, obsBytes)
}
