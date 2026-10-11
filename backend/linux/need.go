//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"time"
)

// need waits until at least count bytes are buffered.
func (x *Driver) need(count int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for x.available() < count {
		if !x.fill(time.Until(deadline)) && (!x.alive || !time.Now().Before(deadline)) {
			return false
		}
	}
	return true
}
