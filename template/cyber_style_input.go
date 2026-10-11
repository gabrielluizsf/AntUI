package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) Input(win *antui.Window, s State, x, y, w, h int, text string, cursor int, caret bool) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	radius := buttonRadius(h, u)

	body, border := cyberPalette.panel, v.edge
	if s.Focused {
		border = v.neon
		glow(cv, x, y, w, h, radius, v.neon, int(float64(20)*v.glow), u)
	} else if s.Hovered {
		body, border = cyberPalette.panelHov, v.neon
	}
	neonFill(cv, x, y, w, h, radius, body, border)

	padding := 8 * u
	cursor = min(max(cursor, 0), len(text))
	cursorPx := canvas.TextWidth(text[:cursor]) * u
	scroll := max(cursorPx-(w-2*padding), 0)

	savedClip := cv.Clip
	cv.SetClip(x+padding, y+u, w-2*padding, h-2*u)

	textY := y + (h-textHeight(u))/2
	colour := v.text
	if text == "" {
		colour = v.dim
	}
	cv.TextScaled(x+padding-scroll, textY, text, colour, u)
	if s.Focused && caret {
		cx := x + padding - scroll + cursorPx
		cv.FillRect(cx, textY, 2*u, textHeight(u), v.neon)
		cv.FillRect(cx, textY+textHeight(u), 2*u, 2*u, canvas.Fade(v.neon, 150))
	}
	cv.Clip = savedClip
}
