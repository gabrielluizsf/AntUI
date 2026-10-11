package template

import "github.com/gabrielluizsf/antui/template/css"

// tableInCell puts a table inside the cell of the row in place. The cell is
// the table's own box, so the row gives it the width of its column and the
// height the table asks for.
func (c *CSS) tableInCell(st css.Style, draw func(*CSS)) {
	if c.measure {
		nested := c.measureTableContainer(st, draw, c.cellRoom())
		cell := c.tableCellSized(st, nested.nw, nested.nh)
		cell.table = nested
		c.row.cells = append(c.row.cells, cell)
		return
	}
	if c.row.next >= len(c.row.cells) {
		return
	}
	cell := c.row.cells[c.row.next]
	c.row.next++
	if cell.table == nil {
		return
	}
	b := c.table
	// A table is as wide as its own columns ask for, not as wide as the cell, so
	// it centres in the cell rather than filling it.
	c.drawTableContainer(cell.table, draw, cell.w, b.originX+cell.x, b.originY+cell.y, cell.h)
}
