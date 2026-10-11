package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// FillTriangle fills a triangle.
func (win *Window) FillTriangle(x0, y0, x1, y1, x2, y2 int, c canvas.Color) {
	win.cv.FillTriangle(x0, y0, x1, y1, x2, y2, c)
}
