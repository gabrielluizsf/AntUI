package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Triangle draws the edges of a triangle.
func (win *Window) Triangle(x0, y0, x1, y1, x2, y2 int, c canvas.Color) {
	win.cv.Triangle(x0, y0, x1, y1, x2, y2, c)
}
