//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"time"
)

// idleWait is how long to sit on the socket with nothing to do: until the
// frame the program asked for is due, and no longer. A program whose frames
// are already late gets no wait at all, which is what keeps a window that
// cannot keep up running flat out instead of stalling between frames.
func (x *Driver) idleWait() time.Duration {
	if x.due.IsZero() {
		return maxIdleWait
	}
	return min(max(time.Until(x.due), 0), maxIdleWait)
}
