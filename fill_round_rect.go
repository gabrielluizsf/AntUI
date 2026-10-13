package antui

import "github.com/gabrielluizsf/antui/canvas"

// FillRoundRect fills a rounded rectangle.
func (win *Window) FillRoundRect(b canvas.Rounded, c canvas.Color) {
	win.cv.FillRoundRect(b, c)
}
