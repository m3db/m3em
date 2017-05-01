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
