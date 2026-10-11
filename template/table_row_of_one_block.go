package template

// tableRowOfOneBlock wraps a block drawn straight in a table in the row it
// belongs to — a row of a single cell, which is what a browser's anonymous
// boxes make of a child that did not ask to be a row. collect runs with the row
// in place, so the block is measured in the silent pass and placed in the draw
// pass like any other cell.
func (c *CSS) tableRowOfOneBlock(collect func(*CSS)) {
	b := c.table
	if b == nil {
		return
	}
	if c.measure {
		r := &tableRow{}
		c.rowWithin(b, r, collect)
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
	c.rowWithin(b, child.row, collect)
}
