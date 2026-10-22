package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) DatePickerBox(p DatePickerBoxParams) {
	win := p.Win
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)

	radius := min(th.Radius*u, p.H/2)
	local := th.Surface
	if p.State.Hovered {
		local = th.SurfaceHover
	}
	border := th.Border
	if p.State.Focused || p.Open {
		border = th.Accent
	}
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, radius), local)
	cv.RoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, radius), border)

	padding := 8 * u
	savedClip := cv.Clip
	cv.SetClip(canvas.Area{X: p.X + padding, Y: p.Y, Width: p.W - 2*padding, Height: p.H})
	cv.TextScaled(canvas.Point{X: float64(p.X + padding), Y: float64(p.Y + (p.H-textHeight(u))/2)}, p.Value, canvas.TextStyle{Color: th.Text, Scale: u})
	cv.Clip = savedClip

	// The caret that says "there is a calendar below", flipped the way the
	// dropdown's is while the picker is open.
	name := IconCaret
	if p.Open {
		name = IconCaretUp
	}
	caretBox := 10 * u
	drawIcon(cv, name,
		p.X+p.W-padding-3*u-caretBox/2, p.Y+(p.H-caretBox)/2, caretBox, caretBox, th.TextMuted)
}
