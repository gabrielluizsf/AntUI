package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// FillCircle fills a disc.
func (win *Window) FillCircle(cx, cy, radius int, c canvas.Color) {
	win.cv.FillCircle(cx, cy, radius, c)
}
