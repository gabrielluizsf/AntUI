package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) SelectOption(win *antui.Window, s State, x, y, w, h int, label string, selected bool) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)

	body := cyberPalette.panel
	if selected || s.Hovered {
		body = cyberPalette.panelHov
	}
	cv.FillRect(canvas.Area{X: x, Y: y, Width: w, Height: h}, body)
	edge := v.edge
	if selected {
		edge = v.neon
	}
	cv.Line(x, y, x+w, y, edge)
	if selected || s.Hovered {
		cv.FillRect(canvas.Area{X: x, Y: y, Width: 3 * u, Height: h}, v.neon)
	}
	colour := v.text
	if selected {
		colour = v.neon
	}
	cv.TextScaled(x+10*u, y+(h-textHeight(u))/2, label, colour, u)
}
