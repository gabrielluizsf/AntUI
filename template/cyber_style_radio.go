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
	if p.State.Focused || p.State.Hovered {
		glow(cv, halo{x: p.X, y: p.Y, w: size, h: size, radius: size / 2, colour: v.neon, strength: int(float64(30) * v.glow), u: u})
	}
	body := cyberPalette.panel
	if p.State.Hovered {
		body = cyberPalette.panelHov
	}
	cv.FillCircle(canvas.Circle{X: cx, Y: cy, R: size / 2}, body)
	border := v.edge
	if p.On || p.State.Hovered {
		border = v.neon
	}
	cv.Circle(canvas.Circle{X: cx, Y: cy, R: size / 2}, border)
	if p.On {
		cv.FillCircle(canvas.Circle{X: cx, Y: cy, R: size/2 - 4*u}, v.neon)
	}
	cv.TextScaled(canvas.Point{X: float64(p.X + size + labelGap*u), Y: float64(p.Y + (size-textHeight(u))/2)}, p.Label, canvas.TextStyle{Color: v.text, Scale: u})
}
