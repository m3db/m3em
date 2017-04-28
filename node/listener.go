package node

import (
	"sync"
	"time"
)

// NewListener creates a new listener
func NewListener(
	onProcessTerminate func(ServiceNode, string),
	onHeartbeatTimeout func(ServiceNode, time.Time),
	onOverwrite func(ServiceNode, string),
) Listener {
	return &listener{
		onProcessTerminate: onProcessTerminate,
		onHeartbeatTimeout: onHeartbeatTimeout,
		onOverwrite:        onOverwrite,
	}
}

type listener struct {
	onProcessTerminate func(ServiceNode, string)
	onHeartbeatTimeout func(ServiceNode, time.Time)
	onOverwrite        func(ServiceNode, string)
}

func (l *listener) OnProcessTerminate(node ServiceNode, desc string) {
	if l.onProcessTerminate != nil {
		l.onProcessTerminate(node, desc)
	}
}

func (l *listener) OnHeartbeatTimeout(node ServiceNode, lastHeartbeatTs time.Time) {
	if l.onHeartbeatTimeout != nil {
		l.onHeartbeatTimeout(node, lastHeartbeatTs)
	}
}

func (l *listener) OnOverwrite(node ServiceNode, desc string) {
	if l.onOverwrite != nil {
		l.onOverwrite(node, desc)
	}
}

type listenerGroup struct {
	sync.Mutex
	node  ServiceNode
	elems map[int]Listener
	token int
}

func newListenerGroup(node ServiceNode) *listenerGroup {
	return &listenerGroup{
		node:  node,
		elems: make(map[int]Listener),
	}
}

func (lg *listenerGroup) add(l Listener) int {
	lg.Lock()
	defer lg.Unlock()
	lg.token++
	lg.elems[lg.token] = l
	return lg.token
}

func (lg *listenerGroup) clear() {
	lg.Lock()
	defer lg.Unlock()
	for i := range lg.elems {
		delete(lg.elems, i)
	}
}

func (lg *listenerGroup) remove(t int) {
	lg.Lock()
	defer lg.Unlock()
	delete(lg.elems, t)
}

func (lg *listenerGroup) notifyTimeout(lastTs time.Time) {
	lg.Lock()
	defer lg.Unlock()
	for _, l := range lg.elems {
		go l.OnHeartbeatTimeout(lg.node, lastTs)
	}
}

func (lg *listenerGroup) notifyTermination(desc string) {
	lg.Lock()
	defer lg.Unlock()
	for _, l := range lg.elems {
		go l.OnProcessTerminate(lg.node, desc)
	}
}

func (lg *listenerGroup) notifyOverwrite(desc string) {
	lg.Lock()
	defer lg.Unlock()
	for _, l := range lg.elems {
		go l.OnOverwrite(lg.node, desc)
	}
}
