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

package fs

import (
	"hash/crc32"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTempFile(t *testing.T, content []byte) *os.File {
	tmpfile, err := ioutil.TempFile("", "example")
	require.NoError(t, err)
	n, err := tmpfile.Write(content)
	require.NoError(t, err)
	require.Equal(t, len(content), n)
	require.NoError(t, tmpfile.Close())
	return tmpfile
}

func TestIterChecksumLargeBuffer(t *testing.T) {
	content := []byte("temporary file content")
	tmpfile := newTempFile(t, content)
	defer os.Remove(tmpfile.Name()) // clean up

	largeBufferSize := 100
	iter, err := NewSizedFileReaderIter(tmpfile.Name(), largeBufferSize)
	require.NoError(t, err)
	var returnedBytes []byte
	numIter := 0
	for iter.Next() {
		returnedBytes = append(returnedBytes, iter.Current()...)
		numIter++
	}
	require.Equal(t, 2, numIter) // once for all the data, once for returning done
	require.NoError(t, iter.Err())
	require.Equal(t, content, returnedBytes)
	require.Equal(t, crc32.ChecksumIEEE(content), iter.Checksum())
	require.Nil(t, iter.(*bufferedFileReaderIter).fileHandle)
}

func TestIterChecksumSmallBuffer(t *testing.T) {
	content := []byte("temporary file content")
	tmpfile := newTempFile(t, content)
	defer os.Remove(tmpfile.Name()) // clean up

	largeBufferSize := 1
	iter, err := NewSizedFileReaderIter(tmpfile.Name(), largeBufferSize)
	require.NoError(t, err)
	numIter := 0
	var returnedBytes []byte
	for iter.Next() {
		returnedBytes = append(returnedBytes, iter.Current()...)
		numIter++
	}
	numExpectedIter := 1 + len(content) // one for each byte, plus one to return done
	require.Equal(t, numExpectedIter, numIter)
	require.NoError(t, iter.Err())
	require.Equal(t, content, returnedBytes)
	require.Equal(t, crc32.ChecksumIEEE(content), iter.Checksum())
	require.Nil(t, iter.(*bufferedFileReaderIter).fileHandle)
}
