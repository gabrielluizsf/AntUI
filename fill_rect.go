package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// FillRect fills a rectangle.
func (win *Window) FillRect(x, y, w, h int, c canvas.Color) {
	win.cv.FillRect(x, y, w, h, c)
}
