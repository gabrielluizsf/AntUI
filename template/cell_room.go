package template

// cellRoom is the width a percentage inside a cell answers to: the space the
// row has for one column of the table, or the room the flow has for the table
// when its columns are not solved yet.
func (c *CSS) cellRoom() int {
	b := c.table
	if b == nil {
		return c.win.Width()
	}
	if len(b.cols) == 0 {
		return max(b.contentW, c.win.Width())
	}
	return b.cols[min(len(c.row.cells), len(b.cols)-1)]
}
