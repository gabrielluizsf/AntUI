package template

import (
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) DatePickerBox(p DatePickerBoxParams) {
	st := cs.styleFor(css.RoleDatePicker, p.State)
	if p.Open {
		st = cs.styleFor(css.RoleDatePicker, State{Focused: true, Hovered: p.State.Hovered})
	}
	cs.paintBox(p.Win, st, p.X, p.Y, p.W, p.H)
	padR := cs.base(st, st.Padding[css.SideRight])
	u := cs.u()
	textW := p.W - padR - 14*u
	cs.drawText(p.Win, st, p.X, p.Y, max(textW, 0), p.H, p.Value)

	// The caret that says "there is a calendar below", flipped the way the
	// dropdown's is while the picker is open.
	name := IconCaret
	if p.Open {
		name = IconCaretUp
	}
	caretBox := 10 * u
	DrawIcon(p.Win.Canvas(), name,
		p.X+p.W-padR-3*u-caretBox/2, p.Y+(p.H-caretBox)/2, caretBox, caretBox, cs.ink(st))
}
