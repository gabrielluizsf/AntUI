package template

import "github.com/gabrielluizsf/antui"

func (cyberStyle) DatePickerBox(win *antui.Window, s State, x, y, w, h int, value string, open bool) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	radius := buttonRadius(h, u)

	body, border := cyberPalette.panel, v.edge
	reach := 0
	if s.Focused || open {
		body, border, reach = cyberPalette.panel, v.neon, 18
	}
	if s.Hovered {
		body = cyberPalette.panelHov
	}
	if reach > 0 {
		glow(cv, halo{x: x, y: y, w: w, h: h, radius: radius, colour: v.neon, strength: int(float64(reach) * v.glow), u: u})
	}
	neonFill(cv, x, y, w, h, radius, body, border)

	padding := 8 * u
	colour := v.text
	if value == "" {
		colour = v.dim
	}
	cv.TextScaled(x+padding, y+(h-textHeight(u))/2, value, colour, u)

	// The caret arrow.
	aw, ah := 6*u, 4*u
	ax := x + w - padding - aw
	ay := y + (h-ah)/2
	if open {
		cv.Line(ax, ay+ah, ax+aw/2, ay, v.neon)
		cv.Line(ax+aw/2, ay, ax+aw, ay+ah, v.neon)
	} else {
		cv.Line(ax, ay, ax+aw/2, ay+ah, v.dim)
		cv.Line(ax+aw/2, ay+ah, ax+aw, ay, v.dim)
	}
}
