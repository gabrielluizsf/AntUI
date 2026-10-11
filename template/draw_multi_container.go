package template

// drawMultiContainer gives the container the box its parent solved for it,
// re-solves it at that size so the columns answer to it, paints its face and
// its rules and draws its children inside.
func (c *CSS) drawMultiContainer(nested *multiNested, draw func(*CSS), pctBase, x, y, w, h int) {
	b := nested.batch
	// The box the parent really gave may be the one the container asked for, in
	// which case it keeps balancing over its columns. A different one — a row
	// taller than its content, a flex item stretched, a height of its own — is a
	// height to cut its boxes against.
	definite := c.multiDeclaredHeight(b.st) > 0 || h != nested.nh
	c.setMultiBox(b, b.st, x, y, w, h, pctBase)
	b.definite = definite
	b.flow()
	if c.style.boxOn(b.st) && visible(b.st) {
		c.style.paintBox(c.win, b.st, x, y, w, b.h)
	}
	c.style.paintColumnRules(c.win, b)
	for pass := range b.passes {
		b.pass = pass
		c.multiWithin(b, draw)
	}
}
