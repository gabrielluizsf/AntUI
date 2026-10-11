package template

import "github.com/gabrielluizsf/antui/template/css"

// natural is the size a widget wants before any CSS width or height says
// otherwise: its text, its padding, its border, at the window's scale.
func (cs *cssStyle) natural(role, label string, st css.Style, u int) (w, h int) {
	return cs.sized(role, st, u, cs.measure(st, label))
}
