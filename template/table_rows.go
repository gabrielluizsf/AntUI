package template

// tableRows gives every row the height of its tallest cell and returns the
// caption's height and the sum of the rows'.
func tableRows(b *tableBatch) (caption, rows int) {
	for _, ch := range b.children {
		if ch.caption != nil {
			caption = max(ch.caption.nh, 0)
		}
	}
	for _, ch := range b.children {
		if ch.row == nil {
			continue
		}
		height := 0
		for _, cell := range ch.row.cells {
			height = max(height, cell.nh)
		}
		ch.row.height = height
		rows += height
	}
	return caption, rows
}
