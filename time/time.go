package time

import (
	"time"
)

// NowFn returns the current time
type NowFn func() time.Time

// ConditionFn specifies a predicate to check
type ConditionFn func() bool

// WaitUntil returns true if the condition specified evaluated to
// true before the timeout, false otherwise
func WaitUntil(fn ConditionFn, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if fn() {
			return true
		}
		time.Sleep(time.Second)
	}
	return false
}
