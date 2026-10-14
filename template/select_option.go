package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) SelectOption(p SelectOptionParams) {
	win := p.Win
	cv := win.Canvas()
	th := win.Theme()

	fill := th.Surface
	border := th.Border
	if p.Selected {
		fill = th.SurfaceHover
		border = th.Accent
	}
	if p.State.Hovered {
		fill = th.AccentHover
		border = th.Accent
	}
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, 0), fill)
	if p.Selected || p.State.Hovered {
		cv.FillRect(canvas.Area{X: p.X, Y: p.Y, Width: 3, Height: p.H}, th.Accent)
	}
	cv.RoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, 0), border)

	u := Scale(win)
	textY := p.Y + (p.H-textHeight(u))/2
	colour := th.Text
	if p.Selected {
		colour = th.TextOnAccent
	}
	cv.TextScaled(p.X+10*u, textY, p.Label, colour, u)
}
