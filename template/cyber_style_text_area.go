package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) TextArea(win *antui.Window, s State, x, y, w, h int, text string, cursor int, caret bool) {
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
	innerW := w - 2*padding
	textW := max(innerW/u, 12)
	cursor = min(max(cursor, 0), len(text))
	lines := textLines(text, textW)

	lineH := textHeight(u)
	visible := max((h-2*padding)/lineH, 1)
	viewTop := textAreaViewTop(cursor, lines, visible)

	savedClip := cv.Clip
	cv.SetClip(x+padding, y+padding, innerW, h-2*padding)

	for lif, l := range lines {
		if lif < viewTop || lif >= viewTop+visible {
			continue
		}
		cy := y + padding + (lif-viewTop)*lineH
		line := text[l[0]:l[1]]
		colour := v.text
		if line == "" {
			colour = v.dim
		}
		cv.TextScaled(x+padding, cy, line, colour, u)
		if s.Focused && caret && cursor >= l[0] && cursor <= l[1] {
			off := cursor - l[0]
			cx := x + padding + textWidth(u, line[:off])
			cv.FillRect(cx, cy, 2*u, lineH, v.neon)
			cv.FillRect(cx, cy+lineH, 2*u, 2*u, canvas.Fade(v.neon, 150))
		}
	}
	cv.Clip = savedClip
}
