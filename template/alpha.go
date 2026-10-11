package template

import "github.com/gabrielluizsf/antui/template/css"

// alpha is the opacity as an alpha fraction of 255.
func (cs *cssStyle) alpha(st css.Style) int {
	if !st.Has("opacity") {
		return 255
	}
	return int(st.Opacity*255 + 0.5)
}
