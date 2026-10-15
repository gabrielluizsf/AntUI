package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) Select(p SelectParams) {
	win := p.Win
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)

	local := th.Surface
	if p.State.Hovered {
		local = th.SurfaceHover
	}
	border := th.Border
	if p.State.Focused || p.Open {
		border = th.Accent
	}
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, th.Radius), local)
	cv.RoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, th.Radius), border)

	padding := 8 * u
	textY := p.Y + (p.H-textHeight(u))/2
	cv.TextScaled(p.X+padding, textY, p.Value, th.Text, u)

	// The arrow that says "there is more below". It is drawn inside a square
	// bigger than the arrow so that the arrow itself comes out at the size the
	// drawing says it is, centred in the square whatever the scale.
	caretBox := 10 * u
	drawIcon(cv, IconCaret,
		p.X+p.W-padding-3*u-caretBox/2, p.Y+(p.H-caretBox)/2, caretBox, caretBox, th.TextMuted)
}
