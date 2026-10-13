package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Select(win *antui.Window, s State, x, y, w, h int, value string, open bool) {
	st := cs.styleFor(css.RoleSelect, s)
	if open {
		st = cs.styleFor(css.RoleSelect, State{Focused: true, Hovered: s.Hovered})
	}
	cs.paintBox(win, st, x, y, w, h)
	padR := cs.base(st, st.Padding[css.SideRight])
	u := cs.u()
	textW := w - padR - 14*u
	cs.drawText(win, st, x, y, max(textW, 0), h, value)

	// The arrow that says "there is more below", drawn as the icon the
	// built-in style draws rather than as two lines, in the box's own ink.
	name := IconCaret
	if open {
		name = IconCaretUp
	}
	caretBox := 10 * u
	DrawIcon(win.Canvas(), name,
		x+w-padR-3*u-caretBox/2, y+(h-caretBox)/2, caretBox, caretBox, cs.ink(st))
}
