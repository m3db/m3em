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
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/m3db/m3em/agent"

	"github.com/stretchr/testify/require"
)

var (
	rngSeed            = int64(123456789)
	defaultUniqueBytes = int64(1024 * 8)
)

func newTempFile(t *testing.T, dir string, content []byte) *os.File {
	tmpfile, err := ioutil.TempFile(dir, "example")
	require.NoError(t, err)
	n, err := tmpfile.Write(content)
	require.NoError(t, err)
	require.Equal(t, len(content), n)
	require.NoError(t, tmpfile.Close())
	return tmpfile
}

func newTestScript(t *testing.T, dir string, scriptNum int, scriptContents []byte) string {
	file, err := ioutil.TempFile(dir, fmt.Sprintf("testscript%d.sh", scriptNum))
	require.NoError(t, err)
	name := file.Name()
	require.NoError(t, file.Chmod(0755))
	numWritten, err := file.Write(scriptContents)
	require.NoError(t, err)
	require.Equal(t, len(scriptContents), numWritten)
	require.NoError(t, file.Close())
	return name
}

func newLargeTempFile(t *testing.T, dir string, numBytes int64) *os.File {
	if numBytes%defaultUniqueBytes != 0 {
		require.FailNow(t, "numBytes must be divisble by %d", defaultUniqueBytes)
	}
	byteStream := newRandByteStream(t, defaultUniqueBytes)
	numIters := numBytes / defaultUniqueBytes
	tmpfile, err := ioutil.TempFile(dir, "example-large-file")
	require.NoError(t, err)
	for i := int64(0); i < numIters; i++ {
		n, err := tmpfile.Write(byteStream)
		require.NoError(t, err)
		require.Equal(t, len(byteStream), n)
	}
	require.NoError(t, tmpfile.Close())
	return tmpfile
}

func newTempDir(t *testing.T) string {
	path, err := ioutil.TempDir("", "integration-test")
	require.NoError(t, err)
	return path
}

func newRandByteStream(t *testing.T, numBytes int64) []byte {
	if mod8 := numBytes % int64(8); mod8 != 0 {
		require.FailNow(t, "numBytes must be divisible by 8")
	}
	var (
		buff = bytes.NewBuffer(nil)
		num  = numBytes / 8
		r    = rand.New(rand.NewSource(rngSeed))
	)
	for i := int64(0); i < num; i++ {
		n := r.Int63()
		err := binary.Write(buff, binary.LittleEndian, n)
		require.NoError(t, err)
	}
	return buff.Bytes()
}

// returns true if agent finished, false otherwise
func waitUntilAgentFinished(a agent.Agent, timeout time.Duration) bool {
	start := time.Now()
	seenRunning := false
	stopped := false
	for !stopped {
		now := time.Now()
		if now.Sub(start) > timeout {
			return stopped
		}
		running := a.Running()
		seenRunning = seenRunning || running
		stopped = seenRunning && !running
		time.Sleep(10 * time.Millisecond)
	}
	return stopped
}
