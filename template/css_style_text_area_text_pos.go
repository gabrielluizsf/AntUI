package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) TextAreaTextPos(win *antui.Window, s State, x, y, w, h int, text string, cursor, mx, my int) int {
	st := cs.styleFor(css.RoleTextArea, s)
	padding := cs.base(st, st.Padding[css.SideLeft])
	textW := max((w-2*padding)/cs.fontScale(st), 12)
	cursor = clampInt(cursor, 0, len(text))
	lines := wrapLines(text, textW, st.WhiteSpace, st.OverflowWrap)
	return textAreaTextPos(text, lines, padding, x, y, h, cursor, cs.lineHeight(st),
		func(prefix string) int { return cs.measure(st, prefix) }, mx, my)
}
