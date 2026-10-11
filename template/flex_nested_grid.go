package template

import "github.com/gabrielluizsf/antui/template/css"

// flexNestedGrid is [CSS.Grid] called from inside a flex: the container becomes
// a flex item, sized and placed by the flex solver like any other.
func (c *CSS) flexNestedGrid(draw func(*CSS)) {
	st := c.style.baseStyle(css.RoleGrid, State{})
	if st.Display == css.DisplayNone {
		return
	}
	parent := c.flex
	if c.measure {
		nested := c.measureGridContainer(st, draw, parent.contentW)
		it := c.flexItemSized(st, nested.nw, nested.nh)
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
