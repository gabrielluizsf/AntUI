package template

import "github.com/gabrielluizsf/antui/template/css"

// flexItemMain records the min and max the item may take along the main axis.
func (c *CSS) flexItemMain(st css.Style, it *flexItem) {
	if c.flex.col() {
		if st.Has("min-height") {
			it.minMain = c.style.length(st, st.MinHeight, c.flex.contentW)
		}
		if st.Has("max-height") {
			it.maxMain = c.style.length(st, st.MaxHeight, c.flex.contentW)
		}
		return
	}
	if st.Has("min-width") {
		it.minMain = c.style.length(st, st.MinWidth, c.flex.contentW)
	}
	if st.Has("max-width") {
		it.maxMain = c.style.length(st, st.MaxWidth, c.flex.contentW)
	}
}
