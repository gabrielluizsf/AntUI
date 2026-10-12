package template

// solveTable gives every cell the width of its column and every row the height
// of its tallest cell, then gives the table the box that makes: a caption above
// the rows, the rows themselves and the table's own border and padding.
//
// x0 and room are the box the table is placed in — the page for a table in the
// flow, the cell, column or flex item for a nested one — and the table centres
// in it the way a block box centres in the flow. pctBase is what a percentage of
// the table itself answers to.
func (c *CSS) solveTable(b *tableBatch, x0, room, pctBase int) {
	winH := c.win.Height()
	left, right, top, bottom := c.tableChrome(b.st, pctBase)
	ml := c.style.length(b.st, b.st.Margin[3], pctBase)
	mr := c.style.length(b.st, b.st.Margin[1], pctBase)
	avail := max(room-ml-mr, 0)

	caption, rows := tableRows(b)
	b.caption = caption

	// The columns ask the cells what they need, and the table's own width — the
	// one it declared, or the one they add up to — is shared back out over them.
	natural := c.tableColumns(b)
	shared := tableShare(natural, c.tableWidth(b, natural, avail, left+right))
	b.cols = shared

	h := caption + rows + top + bottom
	if b.st.Has("height") && gridLengthDefinite(b.st.Height) {
		h = c.clampDim(b.st, c.style.length(b.st, b.st.Height, winH), "height", winH)
	}
	h = max(h, c.clampDim(b.st, 0, "height", winH))

	w := colsWidth(shared) + left + right
	c.setTableBox(b, b.st, x0+ml+(avail-w)/2, b.y, w, h, pctBase)

	placeTableCells(b, shared)
}
