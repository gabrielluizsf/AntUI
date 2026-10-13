package antui

import "github.com/gabrielluizsf/antui/canvas"

// RoundRect draws a rounded outline.
func (win *Window) RoundRect(b canvas.Rounded, c canvas.Color) {
	win.cv.RoundRect(b, c)
}
