package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Button(win *antui.Window, s State, x, y, w, h int, label string) {
	st := cs.styleFor(css.RoleButton, s)
	cs.paintBox(win, st, x, y, w, h)
	cs.drawText(win, st, x, y, w, h, label)
}
