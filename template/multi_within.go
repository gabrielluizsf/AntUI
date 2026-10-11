package template

// multiWithin runs the children of one multi-column container: the callback
// draws them, the silent pass collects them and the draw passes hand them their
// pieces. The batch in place is the container they belong to, so another
// container inside this one is collected as its item, and a container of
// another kind has its own solver put aside while these children are placed.
func (c *CSS) multiWithin(b *multiBatch, draw func(*CSS)) {
	outer, outerMulti, outerFlex, outerGrid, outerTable := c.multi, c.multicoling, c.flexing, c.griding, c.tabling
	c.multi, c.multicoling, c.flexing, c.griding, c.tabling = b, true, false, false, false
	b.next = 0
	draw(c)
	c.multi, c.multicoling, c.flexing, c.griding, c.tabling = outer, outerMulti, outerFlex, outerGrid, outerTable
}
