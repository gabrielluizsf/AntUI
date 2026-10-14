package template

import (
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Label(p LabelParams) {
	st := cs.styleFor(css.RoleLabel, State{})
	cs.run(p.Win.Canvas(), st, p.X, p.Y, p.Text)
}
