package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// Progress paints a read-only progress bar with progress from 0 to 1. It
// reports nothing: a bar has no events.
func (t *uiTemplate) Progress(win *antui.Window, area canvas.Area, progress float32) {
	x, y, w, h := area.X, area.Y, area.Width, area.Height
	t.style.Progress(ProgressParams{
		Win: win,
		X:   x, Y: y, W: w, H: h,
		Value: progress,
	})
}
