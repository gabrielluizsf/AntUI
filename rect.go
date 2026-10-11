package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Rect draws a one-pixel outline.
func (win *Window) Rect(x, y, w, h int, c canvas.Color) { win.cv.Rect(x, y, w, h, c) }
