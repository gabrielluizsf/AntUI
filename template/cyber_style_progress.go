package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) Progress(win *antui.Window, x, y, w, h int, value float32) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)

	value = min(max(value, 0), 1)
	filled := int(float32(w)*value + 0.5)

	cv.FillRoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: w, Height: h}, h/2), cyberPalette.panelAct)
	cv.RoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: w, Height: h}, h/2), v.edge)
	if filled > 0 {
		glow(cv, halo{x: x, y: y, w: filled, h: h, radius: h / 2, colour: v.neon, strength: 5, u: u})
		cv.FillRoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: filled, Height: h}, h/2), v.neon)
	}
}
