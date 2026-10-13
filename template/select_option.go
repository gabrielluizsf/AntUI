package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) SelectOption(win *antui.Window, s State, x, y, w, h int, label string, selected bool) {
	cv := win.Canvas()
	th := win.Theme()

	fill := th.Surface
	border := th.Border
	if selected {
		fill = th.SurfaceHover
		border = th.Accent
	}
	if s.Hovered {
		fill = th.AccentHover
		border = th.Accent
	}
	cv.FillRoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: w, Height: h}, 0), fill)
	if selected || s.Hovered {
		cv.FillRect(canvas.Area{X: x, Y: y, Width: 3, Height: h}, th.Accent)
	}
	cv.RoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: w, Height: h}, 0), border)

	u := Scale(win)
	textY := y + (h-textHeight(u))/2
	colour := th.Text
	if selected {
		colour = th.TextOnAccent
	}
	cv.TextScaled(x+10*u, textY, label, colour, u)
}
