package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (builtinStyle) Progress(p ProgressParams) {
	win := p.Win
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)

	value := min(max(p.Value, 0), 1)
	filled := int(float32(p.W)*value + 0.5)

	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, p.H/2), th.Border)
	if filled > 0 {
		cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: filled, Height: p.H}, p.H/2), th.Accent)
	}
	_ = u
}
