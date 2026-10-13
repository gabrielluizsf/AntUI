package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) TextArea(win *antui.Window, s State, x, y, w, h int, text string, cursor int, caret bool) {
	st := cs.styleFor(css.RoleTextArea, s)
	cs.paintBox(win, st, x, y, w, h)
	u := cs.u()
	padding := cs.base(st, st.Padding[css.SideLeft])
	innerW := w - 2*padding
	textW := max(innerW/cs.fontScale(st), 12)
	cursor = clampInt(cursor, 0, len(text))
	lines := wrapLines(text, textW, st.WhiteSpace, st.OverflowWrap)

	lineH := cs.lineHeight(st)
	visible := max((h-2*padding)/lineH, 1)
	viewTop := 0
	for li, l := range lines {
		if cursor >= l[0] && cursor < l[1] || (li == len(lines)-1 && cursor == l[1]) {
			viewTop = max(li-visible+1, 0)
			break
		}
	}
	viewTop = max(min(viewTop, len(lines)-visible), 0)

	savedClip := win.Canvas().Clip
	win.Canvas().SetClip(x+padding, y+padding, innerW, h-2*padding)

	for lif, l := range lines {
		if lif < viewTop || lif >= viewTop+visible {
			continue
		}
		cy := y + padding + (lif-viewTop)*lineH
		hard := l[1] < len(text) && text[l[1]] == '\n'
		justify := st.TextAlign == 3 && lif < len(lines)-1 && !hard
		// drawTextLine adds the left padding itself; pass the plain x so the
		// line and the caret share the x+padding origin.
		cs.drawTextLine(win, st, x, cy, innerW, lineH, text[l[0]:l[1]], justify)
		if s.Focused && caret && cursor >= l[0] && cursor <= l[1] {
			off := cursor - l[0]
			cx := x + padding + cs.measure(st, text[l[0]:l[0]+off])
			win.Canvas().FillRect(cx, cy, u, lineH, cs.ink(st))
		}
	}
	win.Canvas().Clip = savedClip
}
