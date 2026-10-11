package template

import "github.com/gabrielluizsf/antui/template/css"

// fontScale is how much the style's font-size enlarges the built-in face,
// with 16 reference pixels as the face's natural size.
func (cs *cssStyle) fontScale(st css.Style) int {
	u := cs.u()
	if !st.Has("font-size") || st.FontSize <= 0 {
		return u
	}
	scale := int(float64(st.FontSize*u)/16 + 0.5)
	if scale < 1 {
		scale = 1
	}
	return scale
}
