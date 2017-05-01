package util

import (
	"fmt"
	"sync"
	"time"

	m3dbnode "github.com/m3db/m3em/node/m3db"
	mtime "github.com/m3db/m3em/time"

	xerrors "github.com/m3db/m3x/errors"
)

type nodesWatcher struct {
	sync.Mutex
	pending map[string]m3dbnode.Node
}

// NewM3DBNodesWatcher creates a new M3DBNodeWatcher
func NewM3DBNodesWatcher(nodes []m3dbnode.Node) M3DBNodesWatcher {
	watcher := &nodesWatcher{
		pending: make(map[string]m3dbnode.Node, len(nodes)),
	}
	for _, node := range nodes {
		watcher.addInstance(node)
	}
	return watcher
}

func (nw *nodesWatcher) addInstance(node m3dbnode.Node) {
	nw.Lock()
	nw.pending[node.ID()] = node
	nw.Unlock()
}

func (nw *nodesWatcher) removeInstanceWithLock(id string) {
	delete(nw.pending, id)
}

func (nw *nodesWatcher) Pending() []m3dbnode.Node {
	nw.Lock()
	defer nw.Unlock()

	pending := make([]m3dbnode.Node, 0, len(nw.pending))
	for _, node := range nw.pending {
		pending = append(pending, node)
	}

	return pending
}

func (nw *nodesWatcher) PendingAsError() error {
	nw.Lock()
	defer nw.Unlock()
	var multiErr xerrors.MultiError
	for _, node := range nw.pending {
		multiErr = multiErr.Add(fmt.Errorf("node not bootstrapped: %v", node.ID()))
	}
	return multiErr.FinalError()
}

func (nw *nodesWatcher) WaitUntil(p M3DBNodePredicate, timeout time.Duration) bool {
	nw.Lock()
	defer nw.Unlock()
	var wg sync.WaitGroup
	for id := range nw.pending {
		wg.Add(1)
		m3dbNode := nw.pending[id]
		go func(m3dbNode m3dbnode.Node) {
			defer wg.Done()
			if cond := mtime.WaitUntil(func() bool { return p(m3dbNode) }, timeout); cond {
				nw.removeInstanceWithLock(m3dbNode.ID())
			}
		}(m3dbNode)
	}
	wg.Wait()
	return len(nw.pending) == 0
}
