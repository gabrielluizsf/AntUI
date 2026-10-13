package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) Button(win *antui.Window, s State, x, y, w, h int, label string) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	fill := th.Accent
	switch {
	case s.Pressed:
		fill = canvas.Shade(th.Accent, -0.18)
	case s.Hovered:
		fill = th.AccentHover
	}
	cv.FillRoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: w, Height: h}, th.Radius), fill)

	if s.Focused && !s.Pressed {
		cv.RoundRect(canvas.Round(canvas.Area{X: x - focusRing*u, Y: y - focusRing*u, Width: w + 2*focusRing*u, Height: h + 2*focusRing*u}, th.Radius+focusRing*u), th.Accent)
	}

	textY := y + (h-textHeight(u))/2
	if s.Pressed {
		textY += u
	}
	cv.TextScaled(x+(w-textWidth(u, label))/2, textY, label, th.TextOnAccent, u)
}
