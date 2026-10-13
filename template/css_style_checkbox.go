package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Checkbox(win *antui.Window, s State, x, y int, label string, on bool) {
	st := cs.styleFor(css.RoleCheckbox, s)
	if on {
		st = cs.styleFor(css.RoleCheckbox, State{On: true, Focused: s.Focused, Hovered: s.Hovered})
	}
	u := cs.u()
	box := boxSize * u
	cs.paintBox(win, st, x, y, box, box)
	if on {
		mark := cs.ink(st)
		cv := win.Canvas()
		cv.Line(x+4*u, y+9*u, x+7*u, y+13*u, mark)
		cv.Line(x+5*u, y+9*u, x+8*u, y+13*u, mark)
		cv.Line(x+7*u, y+13*u, x+14*u, y+5*u, mark)
		cv.Line(x+8*u, y+13*u, x+15*u, y+5*u, mark)
	}
	cs.drawText(win, st, x+box+labelGap*u, y, box, box, label)
}
