package template

import "github.com/gabrielluizsf/antui"

func (cyberStyle) Checkbox(win *antui.Window, s State, x, y int, label string, on bool) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	box := boxSize * u
	radius := 4 * u
	if s.Focused || s.Hovered {
		glow(cv, x, y, box, box, radius, v.neon, int(float64(30)*v.glow), u)
	}
	body := cyberPalette.panel
	if s.Hovered {
		body = cyberPalette.panelHov
	}
	cv.FillRoundRect(x, y, box, box, radius, body)
	border := v.edge
	if on || s.Hovered {
		border = v.neon
	}
	cv.RoundRect(x, y, box, box, radius, border)
	if on {
		cv.FillRect(x+4*u, y+9*u, 3*u, 3*u, v.neon)
		cv.FillRect(x+6*u, y+9*u, 8*u, 3*u, v.neon)
	}
	cv.TextScaled(x+box+labelGap*u, y+(box-textHeight(u))/2, label, v.text, u)
}
