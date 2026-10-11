package template

import "github.com/gabrielluizsf/antui/template/css"

// boxOn reports whether the CSS decoration layer has anything to fill a box
// with: a background, a border, outside or inset shadows, an outline, or a
// backdrop filter. A widget whose CSS style has none of these paints its own
// face straight on the window, and the paint path skips the box.
func (cs *cssStyle) boxOn(st css.Style) bool {
	if cs.backgroundOn(st) || st.BorderOn() || cs.outlineOn(st) {
		return true
	}
	return len(st.BoxShadow) > 0 || len(st.BackdropFilters) > 0
}
