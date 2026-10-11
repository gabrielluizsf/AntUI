package template

import "github.com/gabrielluizsf/antui/template/css"

func (c *CSS) flexNestedTable(st css.Style, draw func(*CSS)) {
	parent := c.flex
	if c.measure {
		nested := c.measureTableContainer(st, draw, parent.contentW)
		it := c.flexItemSized(st, nested.nw, nested.nh)
		it.table = nested
		parent.items = append(parent.items, it)
		return
	}
	if parent.next >= len(parent.items) {
		return
	}
	it := &parent.items[parent.next]
	if it.table == nil {
		return
	}
	parent.next++
	c.drawTableContainer(it.table, draw, it.w, parent.originX+it.x, parent.originY+it.y, it.h)
}
