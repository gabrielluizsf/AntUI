package antui

import (
	"time"
)

// Begin starts a frame: it processes system events and refreshes the input
// state. It returns false once the window has been closed.
//
//	for win.Begin() {
//	    win.Clear(canvas.White)
//	    win.End()
//	}
func (win *Window) Begin() bool {
	now := time.Now()
	win.delta = now.Sub(win.frameStart).Seconds()
	if win.delta < 0 || win.delta > 1 {
		win.delta = 1.0 / 60.0
	}
	win.frameStart = now

	win.resetInputFrame()
	win.pollPlatform()
	win.recognize(Now())

	win.cv.ResetClip()
	return !win.shouldClose
}
