package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) TextArea(win *antui.Window, s State, x, y, w, h int, text string, cursor int, caret bool) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	padding := 8 * u

	cv.FillRoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: w, Height: h}, th.Radius), th.Surface)
	border := th.Border
	switch {
	case s.Focused:
		border = th.Accent
	case s.Hovered:
		border = th.TextMuted
	}
	cv.RoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: w, Height: h}, th.Radius), border)

	innerW := w - 2*padding
	textW := max(innerW/u, 12)
	cursor = min(max(cursor, 0), len(text))
	lines := textLines(text, textW)

	lineH := textHeight(u)
	visible := max((h-2*padding)/lineH, 1)
	viewTop := 0
	for li, l := range lines {
		if cursor >= l[0] && cursor < l[1] || (li == len(lines)-1 && cursor == l[1]) {
			viewTop = max(li-visible+1, 0)
			break
		}
	}
	viewTop = max(min(viewTop, len(lines)-visible), 0)

	savedClip := cv.Clip
	cv.SetClip(canvas.Area{X: x + padding, Y: y + padding, Width: innerW, Height: h - 2*padding})

	for lif, l := range lines {
		if lif < viewTop || lif >= viewTop+visible {
			continue
		}
		cy := y + padding + (lif-viewTop)*lineH
		line := text[l[0]:l[1]]
		cv.TextScaled(x+padding, cy, line, th.Text, u)
		if s.Focused && caret && cursor >= l[0] && cursor <= l[1] {
			off := cursor - l[0]
			cx := x + padding + textWidth(u, line[:off])
			cv.FillRect(canvas.Area{X: cx, Y: cy, Width: u, Height: lineH}, th.Text)
		}
	}
	cv.Clip = savedClip
}
