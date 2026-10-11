package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Label(win *antui.Window, x, y int, text string) {
	st := cs.styleFor(css.RoleLabel, State{})
	cs.run(win.Canvas(), st, x, y, text)
}
