package template

import "github.com/gabrielluizsf/antui/template/css"

func (c *CSS) flexNestedMultiCol(st css.Style, draw func(*CSS)) {
	parent := c.flex
	if c.measure {
		nested := c.measureMultiContainer(st, draw, parent.contentW)
		it := c.flexItemSized(st, nested.nw, nested.nh)
		it.columns = nested
		parent.items = append(parent.items, it)
		return
	}
	if parent.next >= len(parent.items) {
		return
	}
	it := &parent.items[parent.next]
	if it.columns == nil {
		return
	}
	parent.next++
	c.drawMultiContainer(it.columns, draw, it.w, parent.originX+it.x, parent.originY+it.y, it.w, it.h)
}
