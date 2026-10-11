package template

import "github.com/gabrielluizsf/antui/template/css"

// measure is how wide text is once the style's transform, weight, slant and
// spacing have been applied.
func (cs *cssStyle) measure(st css.Style, text string) int {
	o := cs.textOpts(st)
	return cs.win.Canvas().TextWidthStyled(transformText(text, st.TextTransform), o)
}
