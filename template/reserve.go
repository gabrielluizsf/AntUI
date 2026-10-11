package template

func (b *gridBatch) reserve(occupied map[[2]int]bool, it *gridItem, col, row, colSpan, rowSpan int) {
	it.colStart, it.colEnd = col, col+colSpan
	it.rowStart, it.rowEnd = row, row+rowSpan
	it.colDefinite, it.rowDefinite = true, true
	b.occupy(occupied, it)
}
