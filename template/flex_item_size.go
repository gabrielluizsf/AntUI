package template

import "github.com/gabrielluizsf/antui/template/css"

// flexItemSize resolves the used width and height from the stylesheet's
// width/height and their min/max bounds, against the flex content box.
func (c *CSS) flexItemSize(st css.Style, nw, nh int) (int, int) {
	b := c.flex
	if st.Has("width") {
		nw = c.style.length(st, st.Width, b.contentW)
	}
	if st.Has("max-width") {
		nw = min(nw, c.style.length(st, st.MaxWidth, b.contentW))
	}
	if st.Has("min-width") {
		nw = max(nw, c.style.length(st, st.MinWidth, b.contentW))
	}
	if st.Has("height") {
		nh = c.style.length(st, st.Height, b.contentW)
	}
	if st.Has("max-height") {
		nh = min(nh, c.style.length(st, st.MaxHeight, b.contentW))
	}
	if st.Has("min-height") {
		nh = max(nh, c.style.length(st, st.MinHeight, b.contentW))
	}
	return nw, nh
}
