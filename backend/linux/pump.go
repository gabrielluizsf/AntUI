//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "github.com/gabrielluizsf/antui/backend"

// Pump handles whatever is waiting on the connection this frame, folding it
// into the window through win.push.
//
// With nothing to say it waits on the socket rather than handing the frame
// straight back, so a program that draws as fast as it can is asleep in the
// kernel between frames instead of asking the server a hundred thousand times
// a second whether anything has happened. The wait ends early on whatever
// arrives, and never runs past the frame the program asked for.
func (x *Driver) Pump(win backend.Face) {
	// Take the whole of what has arrived first, so a frame's worth of input
	// is not spread over the next several frames.
	for x.fill(0) {
	}
	x.dispatch(win)
	if x.alive && x.fill(x.idleWait()) {
		x.dispatch(win)
	}
	if !x.alive {
		win.SetShouldClose()
	}
}
