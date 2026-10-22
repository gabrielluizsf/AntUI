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
	cut := 4 * u
	pts := chamfer(p.X, p.Y, box, box, cut)
	if p.State.Focused {
		chamferGlow(cv, halo{x: p.X, y: p.Y, w: box, h: box, radius: cut, colour: v.neon, strength: int(float64(12) * v.glow), u: u})
	}
	fillClip(cv, pts, cyberPalette.field)
	// Hover lifts the edge one step (acc-dim), checked lights it acc.
	border := v.edge
	if p.State.Hovered {
		border = canvas.Fade(v.neon, 110)
	}
	if p.On {
		border = v.neon
	}
	strokeClip(cv, pts, u, border)
	if p.On {
		drawIconIn(cv, IconCheck, p.X+2*u, p.Y+2*u, 14*u, v.neon)
	}
	cv.TextScaled(canvas.Point{X: float64(p.X + box + labelGap*u), Y: float64(p.Y + (box-textHeight(u))/2)}, p.Label, canvas.TextStyle{Color: v.text, Scale: u})
}
