package environment

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
	elems map[int]Listener
	token int
}

func newListenerGroup() *listenerGroup {
	return &listenerGroup{
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

func (lg *listenerGroup) notifyTimeout(node ServiceNode, lastTs time.Time) {
	lg.Lock()
	defer lg.Unlock()
	for _, l := range lg.elems {
		go l.OnHeartbeatTimeout(node, lastTs)
	}
}

func (lg *listenerGroup) notifyTermination(node ServiceNode, desc string) {
	lg.Lock()
	defer lg.Unlock()
	for _, l := range lg.elems {
		go l.OnProcessTerminate(node, desc)
	}
}

func (lg *listenerGroup) notifyOverwrite(node ServiceNode, desc string) {
	lg.Lock()
	defer lg.Unlock()
	for _, l := range lg.elems {
		go l.OnOverwrite(node, desc)
	}
}
