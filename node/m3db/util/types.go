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

package util

import (
	"time"

	m3dbnode "github.com/m3db/m3em/node/m3db"
)

// M3DBNodePredicate is a predicate on a M3DB ServiceNode
type M3DBNodePredicate func(m3dbnode.Node) bool

// M3DBNodesWatcher makes it easy to monitor observable properties
// of M3DB ServiceNodes
type M3DBNodesWatcher interface {
	// WaitUntil allows you to specify a predicate which must be satisfied
	// on each monitored Node within the timeout provided. It returns a flag
	// indicating if this occurred succesfully
	WaitUntil(p M3DBNodePredicate, timeout time.Duration) bool

	// Pending returns the list of nodes which have not satisfied the
	// predicate satisfied
	Pending() []m3dbnode.Node

	// PendingAsError returns the list of pending nodes wrapped as an
	// error
	PendingAsError() error
}
