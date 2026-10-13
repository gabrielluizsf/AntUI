package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) Checkbox(win *antui.Window, s State, x, y int, label string, on bool) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	box := boxSize * u
	radius := 4 * u
	if s.Focused || s.Hovered {
		glow(cv, halo{x: x, y: y, w: box, h: box, radius: radius, colour: v.neon, strength: int(float64(30) * v.glow), u: u})
	}
	body := cyberPalette.panel
	if s.Hovered {
		body = cyberPalette.panelHov
	}
	cv.FillRoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: box, Height: box}, radius), body)
	border := v.edge
	if on || s.Hovered {
		border = v.neon
	}
	cv.RoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: box, Height: box}, radius), border)
	if on {
		cv.FillRect(canvas.Area{X: x + 4*u, Y: y + 9*u, Width: 3 * u, Height: 3 * u}, v.neon)
		cv.FillRect(canvas.Area{X: x + 6*u, Y: y + 9*u, Width: 8 * u, Height: 3 * u}, v.neon)
	}
	cv.TextScaled(x+box+labelGap*u, y+(box-textHeight(u))/2, label, v.text, u)
}
