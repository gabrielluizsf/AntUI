package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Line draws a line.
func (win *Window) Line(x0, y0, x1, y1 int, c canvas.Color) {
	win.cv.Line(x0, y0, x1, y1, c)
}
