package template

import "github.com/gabrielluizsf/antui/template/css"

func (cs *cssStyle) textHeight(st css.Style) int {
	return cs.textOpts(st).Height()
}
