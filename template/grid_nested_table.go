package template

import "github.com/gabrielluizsf/antui/template/css"

func (c *CSS) gridNestedTable(st css.Style, draw func(*CSS)) {
	parent := c.grid
	if c.measure {
		nested := c.measureTableContainer(st, draw, parent.contentW)
		it := c.gridItemSized("", "", st, nested.nw, nested.nh, nested.nw, nested.nh)
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
