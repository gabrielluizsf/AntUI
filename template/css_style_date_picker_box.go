package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) DatePickerBox(win *antui.Window, s State, x, y, w, h int, value string, open bool) {
	st := cs.styleFor(css.RoleDatePicker, s)
	if open {
		st = cs.styleFor(css.RoleDatePicker, State{Focused: true, Hovered: s.Hovered})
	}
	cs.paintBox(win, st, x, y, w, h)
	padR := cs.base(st, st.Padding[1])
	u := cs.u()
	textW := w - padR - 14*u
	cs.drawText(win, st, x, y, max(textW, 0), h, value)

	aw, ah := 6*u, 4*u
	ax := x + w - padR - aw
	ay := y + (h-ah)/2
	cv := win.Canvas()
	if open {
		cv.Line(ax, ay+ah, ax+aw/2, ay, cs.ink(st))
		cv.Line(ax+aw/2, ay, ax+aw, ay+ah, cs.ink(st))
	} else {
		cv.Line(ax, ay, ax+aw/2, ay+ah, cs.ink(st))
		cv.Line(ax+aw/2, ay+ah, ax+aw, ay, cs.ink(st))
	}
}
