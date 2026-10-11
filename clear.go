package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Clear paints the whole window in one colour.
func (win *Window) Clear(c canvas.Color) { win.cv.Clear(c) }
