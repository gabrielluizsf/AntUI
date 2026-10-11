package template

import "github.com/gabrielluizsf/antui/template/css"

func (cs *cssStyle) padding(st css.Style, side int) int {
	return cs.base(st, st.Padding[side])
}
