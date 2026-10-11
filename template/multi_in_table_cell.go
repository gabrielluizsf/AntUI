package template

import "github.com/gabrielluizsf/antui/template/css"

// multiInTableCell puts a column block inside the cell of the row in place. A
// block drawn straight in the table's callback has no row to land in, so it is
// wrapped in a row of its own first, like a table that founds the same way.
func (c *CSS) multiInTableCell(st css.Style, draw func(*CSS)) {
	if c.row == nil {
		c.tableRowOfOneBlock(func(*CSS) { c.multiInTableCell(st, draw) })
		return
	}
	if c.measure {
		nested := c.measureMultiContainer(st, draw, c.cellRoom())
		cell := c.tableCellSized(st, nested.nw, nested.nh)
		cell.multi = nested
		c.row.cells = append(c.row.cells, cell)
		return
	}
	if c.row.next >= len(c.row.cells) {
		return
	}
	cell := c.row.cells[c.row.next]
	c.row.next++
	if cell.multi == nil {
		return
	}
	b := c.table
	x, y, w, h := b.originX+cell.x, b.originY+cell.y, cell.w, cell.h
	c.drawMultiContainer(cell.multi, draw, w, x, y, w, h)
}
