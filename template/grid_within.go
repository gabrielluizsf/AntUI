package template

// gridWithin runs the children of one grid container: the callback draws them,
// the silent pass collects them and the draw pass hands them their boxes. The
// batch in place is the container they belong to, so a grid nested in this one
// is collected as its item, and a container of another kind has its own solver
// put aside while these children are placed.
func (c *CSS) gridWithin(b *gridBatch, draw func(*CSS)) {
	outer, outerGriding, outerFlexing := c.grid, c.griding, c.flexing
	c.grid, c.griding, c.flexing = b, true, false
	b.next = 0
	draw(c)
	c.grid, c.griding, c.flexing = outer, outerGriding, outerFlexing
}
