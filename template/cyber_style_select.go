package template

func (cyberStyle) Select(p SelectParams) {
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
	textY := p.Y + (p.H-textHeight(u))/2
	colour := v.text
	if p.Value == "" {
		colour = v.dim
	}
	cv.TextScaled(p.X+padding, textY, p.Value, colour, u)

	// The caret arrow.
	aw, ah := 6*u, 4*u
	ax := p.X + p.W - padding - aw
	ay := p.Y + (p.H-ah)/2
	if p.Open {
		cv.Line(ax, ay+ah, ax+aw/2, ay, v.neon)
		cv.Line(ax+aw/2, ay, ax+aw, ay+ah, v.neon)
	} else {
		cv.Line(ax, ay, ax+aw/2, ay+ah, v.dim)
		cv.Line(ax+aw/2, ay+ah, ax+aw, ay, v.dim)
	}
}
