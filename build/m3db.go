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

package build

import (
	"github.com/m3db/m3em/os/fs"
)

type m3dbBuild struct {
	id         string
	sourcePath string
}

func (b *m3dbBuild) ID() string {
	return b.id
}

func (b *m3dbBuild) Iter(bufferSize int) (fs.FileReaderIter, error) {
	return fs.NewSizedFileReaderIter(b.sourcePath, bufferSize)
}

func (b *m3dbBuild) SourcePath() string {
	return b.sourcePath
}

// NewM3DBBuild constructs a new ServiceBuild representing an M3DB build
func NewM3DBBuild(id string, sourcePath string) ServiceBuild {
	return &m3dbBuild{
		id:         id,
		sourcePath: sourcePath,
	}
}

type m3dbConfig struct {
	id    string
	bytes []byte // TODO(prateek): replace with YAML nested structure
}

func (c *m3dbConfig) ID() string {
	return c.id
}

func (c *m3dbConfig) Iter(_ int) (fs.FileReaderIter, error) {
	bytes, err := c.MarshalText()
	if err != nil {
		return nil, err
	}
	return fs.NewBytesReaderIter(bytes), nil
}

func (c *m3dbConfig) MarshalText() ([]byte, error) {
	return c.bytes, nil // TODO(prateek): replace with yaml marshalling
}

func (c *m3dbConfig) UnmarshalText(b []byte) error {
	panic("not implemented")
}

// NewM3DBConfig returns a new M3DB configuration
func NewM3DBConfig(id string, marshalledBytes []byte) ServiceConfiguration {
	return &m3dbConfig{
		id:    id,
		bytes: marshalledBytes,
	}
}
