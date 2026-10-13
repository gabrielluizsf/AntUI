package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) Checkbox(win *antui.Window, s State, x, y int, label string, on bool) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	box := boxSize * u

	fill := th.Surface
	switch {
	case on:
		fill = th.Accent
	case s.Hovered:
		fill = th.SurfaceHover
	}
	border := th.Border
	if on {
		border = th.Accent
	}
	cv.FillRoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: box, Height: box}, 4*u), fill)
	cv.RoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: box, Height: box}, 4*u), border)

	if s.Focused {
		cv.RoundRect(canvas.Round(canvas.Area{X: x - focusRing*u, Y: y - focusRing*u, Width: box + 2*focusRing*u, Height: box + 2*focusRing*u}, focusRingRadius*u), th.Accent)
	}

	if on {
		// The tick is a drawing rather than a few lines so that it is one line of
		// one weight: straight lines drawn to fake a stroke come out thicker on
		// the diagonal than at the corner, which is what a tick drawn four times
		// over used to look like.
		DrawIconIn(cv, IconCheck, x+2*u, y+2*u, 14*u, th.TextOnAccent)
	}
	cv.TextScaled(x+box+labelGap*u, y+(box-textHeight(u))/2, label, th.Text, u)
}
