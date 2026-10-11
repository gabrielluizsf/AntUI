package template

import "github.com/gabrielluizsf/antui"

// Progress paints a read-only progress bar with progress from 0 to 1. It
// reports nothing: a bar has no events.
func (t *uiTemplate) Progress(win *antui.Window, x, y, w, h int, progress float32) {
	t.style.Progress(win, x, y, w, h, progress)
}
