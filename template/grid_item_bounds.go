package template

import "github.com/gabrielluizsf/antui/template/css"

// gridItemBounds records the item's min and max against the grid content box,
// which the solver clamps the used size with.
func (c *CSS) gridItemBounds(st css.Style, it *gridItem, rowDefinite bool) {
	b := c.grid
	if st.Has("min-width") && gridLengthDefinite(st.MinWidth) {
		it.minW = c.style.length(st, st.MinWidth, b.contentW)
	}
	if st.Has("max-width") && gridLengthDefinite(st.MaxWidth) {
		it.maxW = c.style.length(st, st.MaxWidth, b.contentW)
	}
	if st.Has("min-height") && gridLengthDefinite(st.MinHeight) && (!st.MinHeight.IsPct() || rowDefinite) {
		it.minH = c.style.length(st, st.MinHeight, b.contentH)
	}
	if st.Has("max-height") && gridLengthDefinite(st.MaxHeight) && (!st.MaxHeight.IsPct() || rowDefinite) {
		it.maxH = c.style.length(st, st.MaxHeight, b.contentH)
	}
}
