//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

import "time"

// setFrameDeadline tells the driver when the next frame is due, so an idle
// Pump can sleep on the socket until then instead of spinning on it.
func (n *x11Native) setFrameDeadline(due time.Time) { n.d.SetFrameDeadline(due) }
