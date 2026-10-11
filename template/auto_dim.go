package template

import "github.com/gabrielluizsf/antui/template/css"

// autoDim reports whether a background-size length leaves the axis at the
// picture's own size: the auto keyword, the no-value zero Length, or any
// spelling of none.
func autoDim(l css.Length) bool {
	return l.Auto() || l.None() || l == (css.Length{})
}
