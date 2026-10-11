package template

import "github.com/gabrielluizsf/antui"

func (cyberStyle) SelectOption(win *antui.Window, s State, x, y, w, h int, label string, selected bool) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)

	body := cyberPalette.panel
	if selected || s.Hovered {
		body = cyberPalette.panelHov
	}
	cv.FillRect(x, y, w, h, body)
	edge := v.edge
	if selected {
		edge = v.neon
	}
	cv.Line(x, y, x+w, y, edge)
	if selected || s.Hovered {
		cv.FillRect(x, y, 3*u, h, v.neon)
	}
	colour := v.text
	if selected {
		colour = v.neon
	}
	cv.TextScaled(x+10*u, y+(h-textHeight(u))/2, label, colour, u)
}
