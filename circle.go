package antui

import "github.com/gabrielluizsf/antui/canvas"

// Circle draws a circle outline.
func (win *Window) Circle(ci canvas.Circle, c canvas.Color) {
	win.cv.Circle(ci, c)
}
