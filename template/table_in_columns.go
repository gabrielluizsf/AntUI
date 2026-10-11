package template

import "github.com/gabrielluizsf/antui/template/css"

// tableInColumns puts a table inside one column of the multi-column container
// in place. The column is the room the table has, so the table shares its
// columns out over that width, and the flow never cuts it: there is no way to
// know where the table's own rows would be split, so it is laid out whole.
func (c *CSS) tableInColumns(st css.Style, draw func(*CSS)) {
	parent := c.multi
	if c.measure {
		nested := c.measureTableContainer(st, draw, parent.colW)
		it := c.multiItemSized(st, nested.nw, nested.nh)
		it.keep, it.table = true, nested
		parent.items = append(parent.items, it)
		return
	}
	if parent.next >= len(parent.items) {
		return
	}
	it := &parent.items[parent.next]
	parent.next++
	if it.table == nil {
		return
	}
	piece := parent.pieceAt(it)
	x := parent.originX + piece.col*(parent.colW+parent.gap) + it.ml
	y := parent.originY + piece.y
	room := max(parent.colW-it.ml-it.mr, 0)
	// A table is as wide as its columns ask for, not as wide as the column it
	// lands in, so it centres in the column the way any block box is centred.
	c.drawTableContainer(it.table, draw, room, x, y, it.h)
}
