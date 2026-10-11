package template

// TableRow draws one row of a table: the widgets its callback draws are the
// cells of the row, lined up in the columns they belong to. A row is a block of
// its own — the widgets before it and after it are other rows — and the row is
// as tall as the tallest cell in it.
//
// The callback runs in the silent pass to measure the cells and again in the
// draw pass to hand them their boxes, so it must draw the same cells in the same
// order both runs.
func (c *CSS) TableRow(draw func(*CSS)) {
	b := c.table
	if b == nil {
		draw(c)
		return
	}
	if c.measure {
		r := &tableRow{}
		c.rowWithin(b, r, draw)
		b.children = append(b.children, tableChild{row: r})
		return
	}
	if b.rowIdx >= len(b.children) {
		return
	}
	child := b.children[b.rowIdx]
	b.rowIdx++
	if child.row == nil {
		return
	}
	c.rowWithin(b, child.row, draw)
}
