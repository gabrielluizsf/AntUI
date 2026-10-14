package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) TextArea(p TextAreaParams) {
	st := cs.styleFor(css.RoleTextArea, p.State)
	cs.paintBox(p.Win, st, p.X, p.Y, p.W, p.H)
	u := cs.u()
	padding := cs.base(st, st.Padding[css.SideLeft])
	innerW := p.W - 2*padding
	textW := max(innerW/cs.fontScale(st), 12)
	cursor := clampInt(p.Cursor, 0, len(p.Text))
	lines := wrapLines(p.Text, textW, st.WhiteSpace, st.OverflowWrap)

	lineH := cs.lineHeight(st)
	visible := max((p.H-2*padding)/lineH, 1)
	viewTop := 0
	for li, l := range lines {
		if cursor >= l[0] && cursor < l[1] || (li == len(lines)-1 && cursor == l[1]) {
			viewTop = max(li-visible+1, 0)
			break
		}
	}
	viewTop = max(min(viewTop, len(lines)-visible), 0)

	savedClip := p.Win.Canvas().Clip
	p.Win.Canvas().SetClip(canvas.Area{X: p.X + padding, Y: p.Y + padding, Width: innerW, Height: p.H - 2*padding})

	for lif, l := range lines {
		if lif < viewTop || lif >= viewTop+visible {
			continue
		}
		cy := p.Y + padding + (lif-viewTop)*lineH
		hard := l[1] < len(p.Text) && p.Text[l[1]] == '\n'
		justify := st.TextAlign == 3 && lif < len(lines)-1 && !hard
		// drawTextLine adds the left padding itself; pass the plain x so the
		// line and the caret share the x+padding origin.
		cs.drawTextLine(p.Win, st, p.X, cy, innerW, lineH, p.Text[l[0]:l[1]], justify)
		if p.State.Focused && p.Caret && cursor >= l[0] && cursor <= l[1] {
			off := cursor - l[0]
			cx := p.X + padding + cs.measure(st, p.Text[l[0]:l[0]+off])
			p.Win.Canvas().FillRect(canvas.Area{X: cx, Y: cy, Width: u, Height: lineH}, cs.ink(st))
		}
	}
	p.Win.Canvas().Clip = savedClip
}
