package template

// placeTableCells lays every cell out at its column's width and its row's
// height.
func placeTableCells(b *tableBatch, shared []int) {
	rowY := 0
	for _, ch := range b.children {
		if ch.row == nil {
			continue
		}
		x := 0
		for j, cell := range ch.row.cells {
			cell.w = 0
			if j < len(shared) {
				cell.w = shared[j]
			}
			cell.x = x
			cell.h, cell.y = cellPlacement(cell, ch.row.height, rowY)
			x += cell.w
		}
		rowY += ch.row.height
	}
}
