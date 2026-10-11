package template

import "github.com/gabrielluizsf/antui/template/css"

// backgroundOn reports whether a box paints anything behind its border: a
// declared colour, a shorthand that touched the background, or any image
// layer. A box that answers false shows the window's own fill through it.
func (cs *cssStyle) backgroundOn(st css.Style) bool {
	if st.Has("background-color") || st.Has("background") {
		return true
	}
	return len(st.BackgroundImages) > 0
}
