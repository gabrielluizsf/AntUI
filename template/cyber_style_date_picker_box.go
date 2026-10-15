package template

import "github.com/gabrielluizsf/antui/canvas"

func (cyberStyle) DatePickerBox(p DatePickerBoxParams) {
	win := p.Win
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	radius := buttonRadius(p.H, u)

	body, border := cyberPalette.panel, v.edge
	reach := 0
	if p.State.Focused || p.Open {
		body, border, reach = cyberPalette.panel, v.neon, 18
	}
	if p.State.Hovered {
		body = cyberPalette.panelHov
	}
	if reach > 0 {
		glow(cv, halo{x: p.X, y: p.Y, w: p.W, h: p.H, radius: radius, colour: v.neon, strength: int(float64(reach) * v.glow), u: u})
	}
	neonFill(cv, p.X, p.Y, p.W, p.H, radius, body, border)

	padding := 8 * u
	colour := v.text
	if p.Value == "" {
		colour = v.dim
	}
	cv.TextScaled(canvas.Point{X: float64(p.X + padding), Y: float64(p.Y + (p.H-textHeight(u))/2)}, p.Value, canvas.TextStyle{Color: colour, Scale: u})

	// The caret arrow.
	aw, ah := 6*u, 4*u
	ax := p.X + p.W - padding - aw
	ay := p.Y + (p.H-ah)/2
	if p.Open {
		cv.Line(canvas.Point{X: float64(ax), Y: float64(ay + ah)}, canvas.Point{X: float64(ax + aw/2), Y: float64(ay)}, v.neon)
		cv.Line(canvas.Point{X: float64(ax + aw/2), Y: float64(ay)}, canvas.Point{X: float64(ax + aw), Y: float64(ay + ah)}, v.neon)
	} else {
		cv.Line(canvas.Point{X: float64(ax), Y: float64(ay)}, canvas.Point{X: float64(ax + aw/2), Y: float64(ay + ah)}, v.dim)
		cv.Line(canvas.Point{X: float64(ax + aw/2), Y: float64(ay + ah)}, canvas.Point{X: float64(ax + aw), Y: float64(ay)}, v.dim)
	}
}
