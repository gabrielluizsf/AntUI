package template

import (
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Select(p SelectParams) {
	st := cs.styleFor(css.RoleSelect, p.State)
	if p.Open {
		st = cs.styleFor(css.RoleSelect, State{Focused: true, Hovered: p.State.Hovered})
	}
	cs.paintBox(p.Win, st, p.X, p.Y, p.W, p.H)
	padR := cs.base(st, st.Padding[css.SideRight])
	u := cs.u()
	textW := p.W - padR - 14*u
	cs.drawText(p.Win, st, p.X, p.Y, max(textW, 0), p.H, p.Value)

	// The arrow that says "there is more below", drawn as the icon the
	// built-in style draws rather than as two lines, in the box's own ink.
	name := IconCaret
	if p.Open {
		name = IconCaretUp
	}
	caretBox := 10 * u
	DrawIcon(p.Win.Canvas(), name,
		p.X+p.W-padR-3*u-caretBox/2, p.Y+(p.H-caretBox)/2, caretBox, caretBox, cs.ink(st))
}
