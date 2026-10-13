package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) Progress(win *antui.Window, x, y, w, h int, value float32) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)

	value = min(max(value, 0), 1)
	filled := int(float32(w)*value + 0.5)

	cv.FillRoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: w, Height: h}, h/2), th.Border)
	if filled > 0 {
		cv.FillRoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: filled, Height: h}, h/2), th.Accent)
	}
	_ = u
}
