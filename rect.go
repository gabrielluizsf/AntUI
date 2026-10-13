package antui

import "github.com/gabrielluizsf/antui/canvas"

// Rect draws a one-pixel outline.
func (win *Window) Rect(a canvas.Area, c canvas.Color) { win.cv.Rect(a, c) }
