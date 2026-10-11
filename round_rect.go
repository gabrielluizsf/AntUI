package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// RoundRect draws a rounded outline.
func (win *Window) RoundRect(x, y, w, h, r int, c canvas.Color) {
	win.cv.RoundRect(x, y, w, h, r, c)
}
