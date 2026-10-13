//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "time"

// SetFrameDeadline tells the driver when the program wants its next frame,
// which is what an idle Pump waits for. A zero time means the program is not
// pacing itself, and the driver waits its own short moment instead.
func (x *Driver) SetFrameDeadline(due time.Time) { x.due = due }
