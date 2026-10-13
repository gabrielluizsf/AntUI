package antui

import "github.com/gabrielluizsf/antui/canvas"

// FillCircle fills a disc.
func (win *Window) FillCircle(ci canvas.Circle, c canvas.Color) {
	win.cv.FillCircle(ci, c)
}
