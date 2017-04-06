package operator

import (
	"time"
)

var (
	defaultEnabled           = false
	defaultNowFn             = time.Now
	defaultHeartbeatInterval = 10 * time.Second
	defaultCheckInterval     = 2 * time.Second
	defaultHeartbeatTimeout  = 30 * time.Second
)

type heartbeatOpts struct {
	enabled       bool
	nowFn         NowFn
	interval      time.Duration
	checkInterval time.Duration
	timeout       time.Duration
}

// NewHeartbeatOptions returns the default HeartbeatOptions
func NewHeartbeatOptions() HeartbeatOptions {
	return &heartbeatOpts{
		enabled:       defaultEnabled,
		nowFn:         defaultNowFn,
		interval:      defaultHeartbeatInterval,
		checkInterval: defaultCheckInterval,
		timeout:       defaultHeartbeatTimeout,
	}
}

func (ho *heartbeatOpts) SetEnabled(f bool) HeartbeatOptions {
	ho.enabled = f
	return ho
}

func (ho *heartbeatOpts) Enabled() bool {
	return ho.enabled
}

func (ho *heartbeatOpts) SetNowFn(fn NowFn) HeartbeatOptions {
	ho.nowFn = fn
	return ho
}

func (ho *heartbeatOpts) NowFn() NowFn {
	return ho.nowFn
}

func (ho *heartbeatOpts) SetInterval(d time.Duration) HeartbeatOptions {
	ho.interval = d
	return ho
}

func (ho *heartbeatOpts) Interval() time.Duration {
	return ho.interval
}

func (ho *heartbeatOpts) SetCheckInterval(td time.Duration) HeartbeatOptions {
	ho.checkInterval = td
	return ho
}

func (ho *heartbeatOpts) CheckInterval() time.Duration {
	return ho.checkInterval
}

func (ho *heartbeatOpts) SetTimeout(d time.Duration) HeartbeatOptions {
	ho.timeout = d
	return ho
}

func (ho *heartbeatOpts) Timeout() time.Duration {
	return ho.timeout
}
