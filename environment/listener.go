package environment

import (
	"sync"
	"time"
)

// NewListener creates a new listener
func NewListener(
	onProcessTerminate func(string),
	onHeartbeatTimeout func(last time.Time),
	onOverwrite func(string),
) Listener {
	return &listener{
		onProcessTerminate: onProcessTerminate,
		onHeartbeatTimeout: onHeartbeatTimeout,
		onOverwrite:        onOverwrite,
	}
}

type listener struct {
	onProcessTerminate func(string)
	onHeartbeatTimeout func(last time.Time)
	onOverwrite        func(string)
}

func (l *listener) OnProcessTerminate(desc string) {
	if l.onProcessTerminate != nil {
		l.onProcessTerminate(desc)
	}
}

func (l *listener) OnHeartbeatTimeout(lastHeartbeatTs time.Time) {
	if l.onHeartbeatTimeout != nil {
		l.onHeartbeatTimeout(lastHeartbeatTs)
	}
}

func (l *listener) OnOverwrite(desc string) {
	if l.onOverwrite != nil {
		l.onOverwrite(desc)
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

func (lg *listenerGroup) remove(t int) {
	lg.Lock()
	defer lg.Unlock()
	delete(lg.elems, t)
}

func (lg *listenerGroup) notifyTimeout(lastTs time.Time) {
	lg.Lock()
	defer lg.Unlock()
	for _, l := range lg.elems {
		go l.OnHeartbeatTimeout(lastTs)
	}
}

func (lg *listenerGroup) notifyTermination(desc string) {
	lg.Lock()
	defer lg.Unlock()
	for _, l := range lg.elems {
		go l.OnProcessTerminate(desc)
	}
}

func (lg *listenerGroup) notifyOverwrite(desc string) {
	lg.Lock()
	defer lg.Unlock()
	for _, l := range lg.elems {
		go l.OnOverwrite(desc)
	}
}
