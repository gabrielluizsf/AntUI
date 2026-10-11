package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Pixel blends one pixel.
func (win *Window) Pixel(x, y int, c canvas.Color) { win.cv.Pixel(x, y, c) }
