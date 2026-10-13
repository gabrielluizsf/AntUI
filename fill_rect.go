package antui

import "github.com/gabrielluizsf/antui/canvas"

// FillRect fills a rectangle.
func (win *Window) FillRect(a canvas.Area, c canvas.Color) {
	win.cv.FillRect(a, c)
}
