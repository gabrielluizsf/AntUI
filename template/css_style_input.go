package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Input(win *antui.Window, s State, x, y, w, h int, text string, cursor int, caret bool) {
	st := cs.styleFor(css.RoleInput, s)
	cs.paintBox(win, st, x, y, w, h)
	u := cs.u()
	padding := cs.base(st, st.Padding[3])
	savedClip := win.Canvas().Clip
	win.Canvas().SetClip(x+padding, y+1, w-2*padding, h-2)

	cursor = clampInt(cursor, 0, len(text))
	cursorPx := cs.measure(st, text[:cursor])
	scroll := max(cursorPx-(w-2*padding), 0)
	// drawText adds the style's own left padding; pass the pre-scrolled x so
	// the glyphs land at x+padding-scroll — the same origin the caret uses.
	cs.drawText(win, st, x-scroll, y, w, h, text)

	if s.Focused && caret {
		cx := x + padding - scroll + cursorPx
		cy := y + (h-cs.textHeight(st))/2
		win.Canvas().FillRect(cx, cy, u, cs.textHeight(st), cs.ink(st))
	}
	win.Canvas().Clip = savedClip
}
