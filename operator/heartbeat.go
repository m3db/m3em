package operator

import (
	"context"
	"fmt"
	"sync"
	"time"

	hb "github.com/m3db/m3em/generated/proto/heartbeat"

	"github.com/m3db/m3x/instrument"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Heartbeating from agent -> operator: this is to ensure capture of asynchronous
// error conditions, e.g. the child process kicked off by the agent dies, but the
// operator is not informed.
type opHeartbeatServer struct {
	sync.RWMutex
	opts            HeartbeatOptions
	iopts           instrument.Options
	listeners       *listenerGroup
	lastHeartbeat   hb.HeartbeatRequest
	lastHeartbeatTs time.Time
	running         bool
	stopChan        chan bool
	wg              sync.WaitGroup
}

func (h *opHeartbeatServer) Heartbeat(
	ctx context.Context,
	msg *hb.HeartbeatRequest,
) (*hb.HeartbeatResponse, error) {
	nowFn := h.opts.NowFn()

	switch msg.GetCode() {
	case hb.HeartbeatCode_HEALTHY:
		h.updateLastHeartbeat(nowFn(), msg)

	case hb.HeartbeatCode_PROCESS_TERMINATION:
		h.updateLastHeartbeat(nowFn(), msg)
		h.listeners.notifyTermination(msg.GetError())

	case hb.HeartbeatCode_OVERWRITTEN:
		h.listeners.notifyOverwrite(msg.GetError())
		h.stop()

	default:
		return nil, grpc.Errorf(codes.InvalidArgument, "received unknown heartbeat msg: %v", *msg)
	}

	return &hb.HeartbeatResponse{}, nil
}

func newHeartbeater(
	lg *listenerGroup,
	opts HeartbeatOptions,
	iopts instrument.Options,
) *opHeartbeatServer {
	h := &opHeartbeatServer{
		opts:      opts,
		iopts:     iopts,
		listeners: lg,
		stopChan:  make(chan bool, 1),
	}
	return h
}

func (h *opHeartbeatServer) start() error {
	h.Lock()
	defer h.Unlock()
	if h.running {
		return fmt.Errorf("already heartbeating, terminate existing process first")
	}

	h.running = true
	h.wg.Add(1)
	go h.monitorLoop()
	return nil
}

func (h *opHeartbeatServer) lastHeartbeatTime() time.Time {
	h.RLock()
	defer h.RUnlock()
	return h.lastHeartbeatTs
}

func (h *opHeartbeatServer) monitorLoop() {
	defer h.wg.Done()
	var (
		monitor         = true
		checkInterval   = h.opts.CheckInterval()
		timeoutInterval = h.opts.Timeout()
		nowFn           = h.opts.NowFn()
		lastCheck       time.Time
	)

	for monitor {
		select {
		case <-h.stopChan:
			monitor = false
		default:
			last := h.lastHeartbeatTime()
			if last == lastCheck {
				time.Sleep(checkInterval)
				continue
			}
			lastCheck = last
			if !lastCheck.IsZero() && nowFn().Sub(lastCheck) > timeoutInterval {
				h.listeners.notifyTimeout(lastCheck)
			}
			time.Sleep(checkInterval)
		}
	}
}

func (h *opHeartbeatServer) updateLastHeartbeat(ts time.Time, msg *hb.HeartbeatRequest) {
	h.Lock()
	h.lastHeartbeat = *msg
	h.lastHeartbeatTs = ts
	h.Unlock()
}

func (h *opHeartbeatServer) stop() {
	h.Lock()
	if !h.running {
		h.Unlock()
		return
	}

	h.running = false
	h.stopChan <- true
	h.Unlock()
	h.wg.Wait()
}
