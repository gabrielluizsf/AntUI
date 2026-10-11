package template

import "github.com/gabrielluizsf/antui/template/css"

// fontFace is the file the style's font-family list names, or nil when it
// named none — or none this table holds a file for — which is when the canvas
// keeps drawing in the face it already had.
func (cs *cssStyle) fontFace(st css.Style) *css.FontFace {
	if st.FontFamily == "" {
		return nil
	}
	return cs.classes.Font(st.FontFamily, st.FontWeight, st.FontStyle != css.FontStyleNormal)
}
