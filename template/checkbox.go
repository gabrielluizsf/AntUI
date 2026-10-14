package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (builtinStyle) Checkbox(p CheckboxParams) {
	win := p.Win
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	box := boxSize * u

	fill := th.Surface
	switch {
	case p.On:
		fill = th.Accent
	case p.State.Hovered:
		fill = th.SurfaceHover
	}
	border := th.Border
	if p.On {
		border = th.Accent
	}
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: box, Height: box}, 4*u), fill)
	cv.RoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: box, Height: box}, 4*u), border)

	if p.State.Focused {
		cv.RoundRect(canvas.Round(canvas.Area{X: p.X - focusRing*u, Y: p.Y - focusRing*u, Width: box + 2*focusRing*u, Height: box + 2*focusRing*u}, focusRingRadius*u), th.Accent)
	}

	if p.On {
		DrawIconIn(cv, IconCheck, p.X+2*u, p.Y+2*u, 14*u, th.TextOnAccent)
	}
	cv.TextScaled(p.X+box+labelGap*u, p.Y+(box-textHeight(u))/2, p.Label, th.Text, u)
}
