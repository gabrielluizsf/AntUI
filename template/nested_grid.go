package template

import "github.com/gabrielluizsf/antui/template/css"

// nestedGrid is [CSS.Grid] called from inside another grid: the container
// becomes an item of the parent, placed like any other and holding its own
// children.
func (c *CSS) nestedGrid(draw func(*CSS)) {
	st := c.style.baseStyle(css.RoleGrid, State{})
	if st.Display == css.DisplayNone {
		return
	}
	parent := c.grid
	if c.measure {
		nested := c.measureGridContainer(st, draw, parent.contentW)
		it := c.gridItemSized("", "", st, nested.nw, nested.nh, nested.minW, nested.minH)
		it.nested = nested
		parent.items = append(parent.items, it)
		return
	}
	if parent.next >= len(parent.items) {
		return
	}
	it := &parent.items[parent.next]
	if it.nested == nil {
		return
	}
	parent.next++
	c.drawGridContainer(it.nested, draw, parent.contentW, parent.originX+it.x, parent.originY+it.y, it.w, it.h)
}
