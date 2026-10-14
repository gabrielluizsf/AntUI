package template

import (
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) TextAreaTextPos(p TextAreaTextPosParams) int {
	st := cs.styleFor(css.RoleTextArea, p.State)
	padding := cs.base(st, st.Padding[css.SideLeft])
	textW := max((p.W-2*padding)/cs.fontScale(st), 12)
	cursor := clampInt(p.Cursor, 0, len(p.Text))
	lines := wrapLines(p.Text, textW, st.WhiteSpace, st.OverflowWrap)
	return textAreaTextPos(p.Text, lines, padding, p.X, p.Y, p.H, cursor, cs.lineHeight(st),
		func(prefix string) int { return cs.measure(st, prefix) }, p.MouseX, p.MouseY)
}
