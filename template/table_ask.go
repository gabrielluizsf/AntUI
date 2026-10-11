package template

// tableAsk is how wide a table being measured says it wants to be to whatever
// holds it. A width of its own answers for itself; a width in a percentage is a
// fraction of a box the parent has not solved yet — the room the measure pass
// offered is a guess, and in a cell whose columns are not sized it is the whole
// window. A percentage with no definite box to resolve against counts as auto,
// so the answer is the width the columns in the table ask for, which is what
// sizes the column that holds it. The draw pass solves the table again at the
// width the parent really gave, and the percentage answers to that.
func (c *CSS) tableAsk(b *tableBatch, room int) int {
	if !b.st.Has("width") || !b.st.Width.IsPct() {
		return b.w
	}
	left, right, _, _ := c.tableChrome(b.st, room)
	return colsWidth(c.tableColumns(b)) + left + right
}
