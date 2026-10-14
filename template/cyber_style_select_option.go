package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) SelectOption(p SelectOptionParams) {
	win := p.Win
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)

	body := cyberPalette.panel
	if p.Selected || p.State.Hovered {
		body = cyberPalette.panelHov
	}
	cv.FillRect(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, body)
	edge := v.edge
	if p.Selected {
		edge = v.neon
	}
	cv.Line(p.X, p.Y, p.X+p.W, p.Y, edge)
	if p.Selected || p.State.Hovered {
		cv.FillRect(canvas.Area{X: p.X, Y: p.Y, Width: 3 * u, Height: p.H}, v.neon)
	}
	colour := v.text
	if p.Selected {
		colour = v.neon
	}
	cv.TextScaled(p.X+10*u, p.Y+(p.H-textHeight(u))/2, p.Label, colour, u)
}
