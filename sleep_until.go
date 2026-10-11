package antui

import (
	"time"
)

// sleepUntil waits out the rest of the frame. A platform with its own way to
// idle can be given one; the default is the obvious thing.
func (win *Window) sleepUntil(d time.Duration) {
	if win.sleep != nil {
		win.sleep(d)
		return
	}
	time.Sleep(d)
}
