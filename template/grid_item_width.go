package template

import "github.com/gabrielluizsf/antui/template/css"

// gridItemWidth resolves the used width and its track contribution from the
// stylesheet's width, min-width and max-width. A percentage width does not
// contribute to the track's natural size.
func (c *CSS) gridItemWidth(st css.Style, nw, contribW, minContribW int) (int, int, int) {
	b := c.grid
	if st.Has("width") && gridLengthDefinite(st.Width) {
		nw = c.style.length(st, st.Width, b.contentW)
		if !st.Width.IsPct() {
			contribW, minContribW = nw, nw
		}
	}
	if st.Has("max-width") && gridLengthDefinite(st.MaxWidth) {
		limit := c.style.length(st, st.MaxWidth, b.contentW)
		nw = min(nw, limit)
		if !st.MaxWidth.IsPct() {
			contribW = min(contribW, limit)
			minContribW = min(minContribW, limit)
		}
	}
	if st.Has("min-width") && gridLengthDefinite(st.MinWidth) {
		limit := c.style.length(st, st.MinWidth, b.contentW)
		nw = max(nw, limit)
		if !st.MinWidth.IsPct() {
			contribW = max(contribW, limit)
			minContribW = max(minContribW, limit)
		}
	}
	return nw, contribW, minContribW
}
