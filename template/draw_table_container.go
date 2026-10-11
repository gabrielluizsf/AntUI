package template

// drawTableContainer places a nested table in the room its parent solved for it
// — the cell, column or item it was drawn in — solves it again so the columns
// answer to that room, paints its face and lines its rows up inside.
func (c *CSS) drawTableContainer(nested *tableNested, draw func(*CSS), room, x, y, h int) {
	b := nested.batch
	b.y = y
	c.setTableBox(b, b.st, x, y, room, 0, room)
	c.solveTable(b, x, room, room)
	if c.style.boxOn(b.st) && visible(b.st) {
		c.style.paintBox(c.win, b.st, b.x, b.y, b.w, b.h)
	}
	c.tableWithin(b, draw)
}
