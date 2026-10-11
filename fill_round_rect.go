package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// FillRoundRect fills a rounded rectangle.
func (win *Window) FillRoundRect(x, y, w, h, r int, c canvas.Color) {
	win.cv.FillRoundRect(x, y, w, h, r, c)
}
