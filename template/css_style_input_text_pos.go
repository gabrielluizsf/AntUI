package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) InputTextPos(win *antui.Window, s State, x, y, w, h int, text string, cursor, mx, my int) int {
	st := cs.styleFor(css.RoleInput, s)
	padding := cs.base(st, st.Padding[css.SideLeft])
	cursor = clampInt(cursor, 0, len(text))
	cursorPx := cs.measure(st, text[:cursor])
	scroll := max(cursorPx-(w-2*padding), 0)
	return caretAt(func(prefix string) int { return cs.measure(st, prefix) }, text, x+padding-scroll, mx)
}
