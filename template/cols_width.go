package template

// colsWidth is what the columns add up to.
func colsWidth(cols []int) int {
	total := 0
	for _, w := range cols {
		total += w
	}
	return total
}
