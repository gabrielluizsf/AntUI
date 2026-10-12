package template

// gridCellSize is the span the item covers: the tracks it crosses, plus the
// gaps inside the span.
func (b *gridBatch) gridCellSize(it *gridItem, columns, rows []int) (int, int) {
	cellW, cellH := 0, 0
	for col := it.colStart; col < it.colEnd && col < len(columns); col++ {
		cellW += columns[col]
	}
	for row := it.rowStart; row < it.rowEnd && row < len(rows); row++ {
		cellH += rows[row]
	}
	cellW += max(it.colEnd-it.colStart-1, 0) * b.colGap
	cellH += max(it.rowEnd-it.rowStart-1, 0) * b.rowGap
	return cellW, cellH
}
