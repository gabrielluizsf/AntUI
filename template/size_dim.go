package template

import "github.com/gabrielluizsf/antui/template/css"

// sizeDim measures one background-size axis against its positioning area,
// falling back onto the picture's own width or height.
func (cs *cssStyle) sizeDim(st css.Style, l css.Length, base int, natural int) int {
	if autoDim(l) {
		return natural
	}
	if l.IsPct() {
		return l.Resolve(css.Units{Width: base})
	}
	return l.Resolve(css.Units{Width: base}) * cs.u()
}
