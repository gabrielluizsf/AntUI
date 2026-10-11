package template

// tableColumns is the natural width of each column: the widest cell that lands
// in it, over every row of the table.
func (c *CSS) tableColumns(b *tableBatch) []int {
	var cols []int
	for _, ch := range b.children {
		if ch.row == nil {
			continue
		}
		for j, cell := range ch.row.cells {
			for len(cols) <= j {
				cols = append(cols, 0)
			}
			cols[j] = max(cols[j], cell.nw)
		}
	}
	return cols
}
