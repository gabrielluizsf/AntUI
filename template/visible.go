package template

import "github.com/gabrielluizsf/antui/template/css"

// visible reports whether a widget painted by the stylesheet should draw and
// hear input. visibility:hidden keeps its flow slot but is skipped here.
func visible(st css.Style) bool {
	return st.Visibility == css.VisibilityVisible
}
