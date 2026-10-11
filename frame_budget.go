package antui

import (
	"time"
)

// frameBudget is how long one frame may take. The default is the display's own
// refresh rate, so the loop is paced by the screen rather than by how fast the
// machine can go; a rate given to [Window.SetFPS] is a ceiling on top of that.
// A display that will not say is assumed to be the usual sixty.
func (win *Window) frameBudget() time.Duration {
	fps := win.targetFPS
	if fps <= 0 {
		if win.refreshHz == 0 && win.native != nil {
			win.refreshHz = win.native.displayRefresh()
		}
		if fps = win.refreshHz; fps <= 0 {
			fps = fallbackFPS
		}
	}
	return time.Duration(float64(time.Second) / float64(fps))
}
