package template

// drawGridContainer gives the container the box its parent solved for it,
// re-solves it at that size so the tracks answer to it, paints its face and
// draws its children inside.
func (c *CSS) drawGridContainer(nested *gridNested, draw func(*CSS), pctBase, x, y, w, h int) {
	b := nested.batch
	c.setGridBox(b, b.st, x, y, w, h, pctBase)
	b.solve()
	if c.style.boxOn(b.st) && visible(b.st) {
		c.style.paintBox(c.win, b.st, x, y, w, h)
	}
	c.gridWithin(b, draw)
}
