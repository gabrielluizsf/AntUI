package antui

import "github.com/gabrielluizsf/antui/canvas"

// Line draws a straight line.
func (win *Window) Line(a, b canvas.Point, c canvas.Color) {
	win.cv.Line(a, b, c)
}
