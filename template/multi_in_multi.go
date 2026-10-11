package template

import "github.com/gabrielluizsf/antui/template/css"

func (c *CSS) multiInMulti(st css.Style, draw func(*CSS)) {
	parent := c.multi
	if c.measure {
		nested := c.measureMultiContainer(st, draw, parent.contentW)
		it := c.multiItemSized(st, nested.nw, nested.nh)
		it.keep, it.nested = true, nested
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
	piece := parent.pieceAt(it)
	x := parent.originX + it.ml
	y := parent.originY + piece.y
	c.drawMultiContainer(it.nested, draw, parent.colW, x, y, max(parent.colW-it.ml-it.mr, 0), it.h)
}
