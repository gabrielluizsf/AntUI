package template

import "github.com/gabrielluizsf/antui/template/css"

// gridItemInsets records where an out-of-flow item asks to be pinned.
func (c *CSS) gridItemInsets(st css.Style, it *gridItem) {
	b := c.grid
	if !st.OutOfFlow() {
		return
	}
	it.out = true
	if st.Has("left") {
		it.inLeft = c.style.length(st, st.Left, b.contentW)
	}
	if st.Has("right") {
		it.inRight = c.style.length(st, st.Right, b.contentW)
	}
	if st.Has("top") {
		it.inTop = c.style.length(st, st.Top, b.contentW)
	}
	if st.Has("bottom") {
		it.inBottom = c.style.length(st, st.Bottom, b.contentW)
	}
}
