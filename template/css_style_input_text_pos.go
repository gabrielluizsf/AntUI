package template

import (
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) InputTextPos(p InputTextPosParams) int {
	st := cs.styleFor(css.RoleInput, p.State)
	padding := cs.base(st, st.Padding[css.SideLeft])
	cursor := clampInt(p.Cursor, 0, len(p.Text))
	cursorPx := cs.measure(st, p.Text[:cursor])
	scroll := max(cursorPx-(p.W-2*padding), 0)
	return caretAt(func(prefix string) int { return cs.measure(st, prefix) }, p.Text, p.X+padding-scroll, p.MouseX)
}
