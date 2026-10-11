package template

// tableWithin runs the children of one table: the callback draws them, the
// silent pass collects the rows and the draw pass hands the cells their boxes.
// The batch in place is the table they belong to, so a table inside this one is
// collected as a cell of its own. The pass the page is in carries through, as it
// does in multiWithin: a nested table reached by the silent pass collects its own
// rows there, and one reached by a draw pass places them.
func (c *CSS) tableWithin(b *tableBatch, draw func(*CSS)) {
	outer, outerTable, outerRow, outerMulti := c.table, c.tabling, c.row, c.multicoling
	outerFlex, outerGrid := c.flexing, c.griding
	c.table, c.tabling, c.row, c.multicoling = b, true, nil, false
	c.flexing, c.griding = false, false
	b.rowIdx = 0
	draw(c)
	c.table, c.tabling, c.row, c.multicoling = outer, outerTable, outerRow, outerMulti
	c.flexing, c.griding = outerFlex, outerGrid
}
