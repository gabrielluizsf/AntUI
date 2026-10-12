package template

import "github.com/gabrielluizsf/antui/template/css"

// gridItemHeight resolves the used height and its track contribution from the
// stylesheet's height, min-height and max-height. A percentage height only
// counts when the row axis is definite, since there is nothing to measure a
// percentage against in a growing grid.
func (c *CSS) gridItemHeight(st css.Style, nh, contribH, minContribH int, rowDefinite bool) (int, int, int) {
	b := c.grid
	if st.Has("height") && gridLengthDefinite(st.Height) && (!st.Height.IsPct() || rowDefinite) {
		nh = c.style.length(st, st.Height, b.contentH)
		if !st.Height.IsPct() {
			contribH, minContribH = nh, nh
		}
	}
	if st.Has("max-height") && gridLengthDefinite(st.MaxHeight) && (!st.MaxHeight.IsPct() || rowDefinite) {
		limit := c.style.length(st, st.MaxHeight, b.contentH)
		nh = min(nh, limit)
		if !st.MaxHeight.IsPct() {
			contribH = min(contribH, limit)
			minContribH = min(minContribH, limit)
		}
	}
	if st.Has("min-height") && gridLengthDefinite(st.MinHeight) && (!st.MinHeight.IsPct() || rowDefinite) {
		limit := c.style.length(st, st.MinHeight, b.contentH)
		nh = max(nh, limit)
		if !st.MinHeight.IsPct() {
			contribH = max(contribH, limit)
			minContribH = max(minContribH, limit)
		}
	}
	return nh, contribH, minContribH
}
