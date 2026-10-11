package template

import "github.com/gabrielluizsf/antui/template/css"

// perSideBorder reports whether the four sides cannot be drawn as one ring:
// a non-solid style, a missing side, or sides of different width or colour.
func (cs *cssStyle) perSideBorder(st css.Style) bool {
	first := -1
	for i := 0; i < 4; i++ {
		if st.BoxStyle[i] == css.BorderNone {
			return true
		}
		if st.BoxStyle[i] > css.BorderSolid {
			return true
		}
		if first < 0 {
			first = i
			continue
		}
		if st.BorderWidth[i] != st.BorderWidth[first] || st.BoxColor[i] != st.BoxColor[first] {
			return true
		}
	}
	return false
}
