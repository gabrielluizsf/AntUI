package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) SelectOption(p SelectOptionParams) {
	win := p.Win
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)

	// The site's list rows: hover tints, selected lights the accent bar.
	body := cyberPalette.panel
	if p.Selected {
		body = canvas.Fade(v.neon, 30)
	} else if p.State.Hovered {
		body = canvas.Fade(v.neon, 20)
	}
	cv.FillRect(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, body)
	cv.Line(canvas.Point{X: float64(p.X), Y: float64(p.Y)}, canvas.Point{X: float64(p.X + p.W), Y: float64(p.Y)}, v.edge)
	if p.Selected {
		cv.FillRect(canvas.Area{X: p.X, Y: p.Y, Width: 3 * u, Height: p.H}, v.neon)
	}
	cv.TextScaled(canvas.Point{X: float64(p.X + 10*u), Y: float64(p.Y + (p.H-textHeight(u))/2)}, p.Label, canvas.TextStyle{Color: v.text, Scale: u})
}
