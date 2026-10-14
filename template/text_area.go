package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) TextArea(p TextAreaParams) {
	win := p.Win
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	padding := 8 * u

	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, th.Radius), th.Surface)
	border := th.Border
	switch {
	case p.State.Focused:
		border = th.Accent
	case p.State.Hovered:
		border = th.TextMuted
	}
	cv.RoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, th.Radius), border)

	innerW := p.W - 2*padding
	textW := max(innerW/u, 12)
	cursor := min(max(p.Cursor, 0), len(p.Text))
	lines := textLines(p.Text, textW)

	lineH := textHeight(u)
	visible := max((p.H-2*padding)/lineH, 1)
	viewTop := 0
	for li, l := range lines {
		if cursor >= l[0] && cursor < l[1] || (li == len(lines)-1 && cursor == l[1]) {
			viewTop = max(li-visible+1, 0)
			break
		}
	}
	viewTop = max(min(viewTop, len(lines)-visible), 0)

	savedClip := cv.Clip
	cv.SetClip(canvas.Area{X: p.X + padding, Y: p.Y + padding, Width: innerW, Height: p.H - 2*padding})

	for lif, l := range lines {
		if lif < viewTop || lif >= viewTop+visible {
			continue
		}
		cy := p.Y + padding + (lif-viewTop)*lineH
		line := p.Text[l[0]:l[1]]
		cv.TextScaled(p.X+padding, cy, line, th.Text, u)
		if p.State.Focused && p.Caret && cursor >= l[0] && cursor <= l[1] {
			off := cursor - l[0]
			cx := p.X + padding + textWidth(u, line[:off])
			cv.FillRect(canvas.Area{X: cx, Y: cy, Width: u, Height: lineH}, th.Text)
		}
	}
	cv.Clip = savedClip
}
