package template

// rowWithin runs the cells of one row: the callback draws them and the row
// collects them. A row that is not the one in place is put back afterwards, so
// a table inside a cell cannot leave the outer row collecting.
func (c *CSS) rowWithin(b *tableBatch, r *tableRow, draw func(*CSS)) {
	outer := c.row
	c.row = r
	r.next = 0
	draw(c)
	c.row = outer
}
