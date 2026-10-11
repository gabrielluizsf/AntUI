package template

import "github.com/gabrielluizsf/antui/template/css"

// interactive reports whether a widget takes pointer input. pointer-events:
// none still paints, so the caller falls back to a paint-only path.
func interactive(st css.Style) bool {
	return st.PointerEvents != css.PointerEventsNone
}
