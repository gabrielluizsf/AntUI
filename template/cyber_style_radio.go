package template

import "github.com/gabrielluizsf/antui"

func (cyberStyle) Radio(win *antui.Window, s State, x, y int, label string, on bool) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	size := 18 * u
	cx, cy := x+size/2, y+size/2
	if s.Focused || s.Hovered {
		glow(cv, x, y, size, size, size/2, v.neon, int(float64(30)*v.glow), u)
	}
	body := cyberPalette.panel
	if s.Hovered {
		body = cyberPalette.panelHov
	}
	cv.FillCircle(cx, cy, size/2, body)
	border := v.edge
	if on || s.Hovered {
		border = v.neon
	}
	cv.Circle(cx, cy, size/2, border)
	if on {
		cv.FillCircle(cx, cy, size/2-4*u, v.neon)
	}
	cv.TextScaled(x+size+8*u, y+(size-textHeight(u))/2, label, v.text, u)
}
