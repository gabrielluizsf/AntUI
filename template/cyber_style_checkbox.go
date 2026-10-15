package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) Checkbox(p CheckboxParams) {
	win := p.Win
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	box := boxSize * u
	radius := 4 * u
	if p.State.Focused || p.State.Hovered {
		glow(cv, halo{x: p.X, y: p.Y, w: box, h: box, radius: radius, colour: v.neon, strength: int(float64(30) * v.glow), u: u})
	}
	body := cyberPalette.panel
	if p.State.Hovered {
		body = cyberPalette.panelHov
	}
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: box, Height: box}, radius), body)
	border := v.edge
	if p.On || p.State.Hovered {
		border = v.neon
	}
	cv.RoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: box, Height: box}, radius), border)
	if p.On {
		cv.FillRect(canvas.Area{X: p.X + 4*u, Y: p.Y + 9*u, Width: 3 * u, Height: 3 * u}, v.neon)
		cv.FillRect(canvas.Area{X: p.X + 6*u, Y: p.Y + 9*u, Width: 8 * u, Height: 3 * u}, v.neon)
	}
	cv.TextScaled(canvas.Point{X: float64(p.X + box + labelGap*u), Y: float64(p.Y + (box-textHeight(u))/2)}, p.Label, canvas.TextStyle{Color: v.text, Scale: u})
}
