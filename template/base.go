package template

import "github.com/gabrielluizsf/antui/template/css"

// base measures a length against the window width; used for padding and
// margins.
func (cs *cssStyle) base(st css.Style, l css.Length) int {
	return cs.length(st, l, cs.win.Width())
}
