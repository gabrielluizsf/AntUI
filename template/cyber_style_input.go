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
		glow(cv, halo{x: x, y: y, w: w, h: h, radius: radius, colour: v.neon, strength: int(float64(20) * v.glow), u: u})
	} else if s.Hovered {
		body, border = cyberPalette.panelHov, v.neon
	}
	neonFill(cv, x, y, w, h, radius, body, border)

	padding := 8 * u
	cursor = min(max(cursor, 0), len(text))
	cursorPx := canvas.TextWidth(text[:cursor]) * u
	scroll := max(cursorPx-(w-2*padding), 0)

	savedClip := cv.Clip
	cv.SetClip(canvas.Area{X: x + padding, Y: y + u, Width: w - 2*padding, Height: h - 2*u})

	textY := y + (h-textHeight(u))/2
	colour := v.text
	if text == "" {
		colour = v.dim
	}
	cv.TextScaled(x+padding-scroll, textY, text, colour, u)
	if s.Focused && caret {
		cx := x + padding - scroll + cursorPx
		cv.FillRect(canvas.Area{X: cx, Y: textY, Width: 2 * u, Height: textHeight(u)}, v.neon)
		cv.FillRect(canvas.Area{X: cx, Y: textY + textHeight(u), Width: 2 * u, Height: 2 * u}, canvas.Fade(v.neon, 150))
	}
	cv.Clip = savedClip
}
