package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) Radio(p RadioParams) {
	win := p.Win
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	size := boxSize * u
	cx, cy := p.X+size/2, p.Y+size/2
	if p.State.Focused {
		glow(cv, halo{x: p.X, y: p.Y, w: size, h: size, radius: size / 2, colour: v.neon, strength: int(float64(12) * v.glow), u: u})
	}
	cv.FillCircle(canvas.Circle{X: cx, Y: cy, R: size / 2}, cyberPalette.field)
	// Hover lifts the edge one step (acc-dim), selected lights it acc.
	border := v.edge
	if p.State.Hovered {
		border = canvas.Fade(v.neon, 110)
	}
	if p.On {
		border = v.neon
	}
	cv.Circle(canvas.Circle{X: cx, Y: cy, R: size / 2}, border)
	if p.On {
		cv.FillCircle(canvas.Circle{X: cx, Y: cy, R: size/2 - 4*u}, v.neon)
	}
	cv.TextScaled(canvas.Point{X: float64(p.X + size + labelGap*u), Y: float64(p.Y + (size-textHeight(u))/2)}, p.Label, canvas.TextStyle{Color: v.text, Scale: u})
}
