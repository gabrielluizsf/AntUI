package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) Input(win *antui.Window, s State, x, y, w, h int, text string, cursor int, caret bool) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	padding := 8 * u

	cv.FillRoundRect(x, y, w, h, th.Radius, th.Surface)
	border := th.Border
	switch {
	case s.Focused:
		border = th.Accent
	case s.Hovered:
		border = th.TextMuted
	}
	cv.RoundRect(x, y, w, h, th.Radius, border)

	cursor = min(max(cursor, 0), len(text))
	cursorPx := canvas.TextWidth(text[:cursor]) * u
	scroll := max(cursorPx-(w-2*padding), 0)

	savedClip := cv.Clip
	cv.SetClip(x+padding, y+u, w-2*padding, h-2*u)

	textY := y + (h-textHeight(u))/2
	cv.TextScaled(x+padding-scroll, textY, text, th.Text, u)

	if s.Focused && caret {
		cv.FillRect(x+padding-scroll+cursorPx, textY-u, u, textHeight(u)+2*u, th.Text)
	}
	cv.Clip = savedClip
}
