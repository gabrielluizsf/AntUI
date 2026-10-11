package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Circle draws a circle outline.
func (win *Window) Circle(cx, cy, radius int, c canvas.Color) {
	win.cv.Circle(cx, cy, radius, c)
}
