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
	area := canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: box, Height: box}, radius)
	if p.State.Focused {
		glow(cv, halo{x: p.X, y: p.Y, w: box, h: box, radius: radius, colour: v.neon, strength: int(float64(12) * v.glow), u: u})
	}
	cv.FillRoundRect(area, cyberPalette.field)
	// Hover lifts the edge one step (acc-dim), checked lights it acc.
	border := v.edge
	if p.State.Hovered {
		border = canvas.Fade(v.neon, 110)
	}
	if p.On {
		border = v.neon
	}
	cv.RoundRect(area, border)
	if p.On {
		cv.FillRect(canvas.Area{X: p.X + 4*u, Y: p.Y + 9*u, Width: 3 * u, Height: 3 * u}, v.neon)
		cv.FillRect(canvas.Area{X: p.X + 6*u, Y: p.Y + 9*u, Width: 8 * u, Height: 3 * u}, v.neon)
	}
	cv.TextScaled(canvas.Point{X: float64(p.X + box + labelGap*u), Y: float64(p.Y + (box-textHeight(u))/2)}, p.Label, canvas.TextStyle{Color: v.text, Scale: u})
}
