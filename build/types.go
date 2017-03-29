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

// ServiceBuild represents a single build of a service, available locally.
type ServiceBuild interface {
	// ID returns a string identifier for the build
	ID() string

	// SourcePath returns a local path where the service build is available
	SourcePath() string
}

// ServiceConfiguration represents the configuration required to operate a service
type ServiceConfiguration interface {
	// ID returns a string identifier for the service configuration
	ID() string

	// MarshalText returns a UTF-8 serialized representation of the service configuration
	MarshalText() ([]byte, error)

	// UnmarshalText deserializes the provided byte stream
	UnmarshalText([]byte) error
}
