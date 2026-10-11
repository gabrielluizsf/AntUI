package template

func (b *gridBatch) fits(occupied map[[2]int]bool, col, row, colSpan, rowSpan int) bool {
	if col < 0 || row < 0 || colSpan < 1 || rowSpan < 1 {
		return false
	}
	for r := row; r < row+rowSpan; r++ {
		for c := col; c < col+colSpan; c++ {
			if occupied[[2]int{c, r}] {
				return false
			}
		}
	}
	return true
}
