package template

import (
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Button(p ButtonParams) {
	st := cs.styleFor(css.RoleButton, p.State)
	cs.paintBox(p.Win, st, p.X, p.Y, p.W, p.H)
	cs.drawText(p.Win, st, p.X, p.Y, p.W, p.H, p.Label)
}
