package template

import "github.com/gabrielluizsf/antui/template/css"

// multiItemBox is the width and height a child takes in a column colW wide: the
// size it asked for, capped by max and lifted by min. The flow calls it for
// every solve, so a width in a percentage answers to the column the child
// really ended up in.
func (c *CSS) multiItemBox(st css.Style, nw, nh, colW int) (w, h int) {
	if st.Has("width") {
		nw = c.style.length(st, st.Width, colW)
	}
	if st.Has("max-width") {
		nw = min(nw, c.style.length(st, st.MaxWidth, colW))
	}
	if st.Has("min-width") {
		nw = max(nw, c.style.length(st, st.MinWidth, colW))
	}
	if st.Has("height") {
		nh = c.style.length(st, st.Height, colW)
	}
	if st.Has("max-height") {
		nh = min(nh, c.style.length(st, st.MaxHeight, colW))
	}
	if st.Has("min-height") {
		nh = max(nh, c.style.length(st, st.MinHeight, colW))
	}
	return nw, nh
}
