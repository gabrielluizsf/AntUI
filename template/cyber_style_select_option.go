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
	cv.Line(canvas.Point{X: float64(p.X), Y: float64(p.Y)}, canvas.Point{X: float64(p.X + p.W), Y: float64(p.Y)}, edge)
	if p.Selected || p.State.Hovered {
		cv.FillRect(canvas.Area{X: p.X, Y: p.Y, Width: 3 * u, Height: p.H}, v.neon)
	}
	colour := v.text
	if p.Selected {
		colour = v.neon
	}
	cv.TextScaled(canvas.Point{X: float64(p.X + 10*u), Y: float64(p.Y + (p.H-textHeight(u))/2)}, p.Label, canvas.TextStyle{Color: colour, Scale: u})
}
