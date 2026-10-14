package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) Progress(p ProgressParams) {
	win := p.Win
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)

	value := min(max(p.Value, 0), 1)
	filled := int(float32(p.W)*value + 0.5)

	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, p.H/2), cyberPalette.panelAct)
	cv.RoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, p.H/2), v.edge)
	if filled > 0 {
		glow(cv, halo{x: p.X, y: p.Y, w: filled, h: p.H, radius: p.H / 2, colour: v.neon, strength: 5, u: u})
		cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: filled, Height: p.H}, p.H/2), v.neon)
	}
}
