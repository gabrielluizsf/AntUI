package template

func (b *gridBatch) occupy(occupied map[[2]int]bool, it *gridItem) {
	for row := it.rowStart; row < it.rowEnd; row++ {
		for col := it.colStart; col < it.colEnd; col++ {
			occupied[[2]int{col, row}] = true
		}
	}
}
