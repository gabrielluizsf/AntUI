package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) Button(p ButtonParams) {
	cv := p.Win.Canvas()
	th := p.Win.Theme()
	u := Scale(p.Win)
	fill := th.Accent
	switch {
	case p.State.Pressed:
		fill = canvas.Shade(th.Accent, -0.18)
	case p.State.Hovered:
		fill = th.AccentHover
	}
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, th.Radius), fill)

	if p.State.Focused && !p.State.Pressed {
		cv.RoundRect(canvas.Round(canvas.Area{X: p.X - focusRing*u, Y: p.Y - focusRing*u, Width: p.W + 2*focusRing*u, Height: p.H + 2*focusRing*u}, th.Radius+focusRing*u), th.Accent)
	}

	textY := p.Y + (p.H-textHeight(u))/2
	if p.State.Pressed {
		textY += u
	}
	cv.TextScaled(canvas.Point{X: float64(p.X + (p.W-textWidth(u, p.Label))/2), Y: float64(textY)}, p.Label, canvas.TextStyle{Color: th.TextOnAccent, Scale: u})
}
