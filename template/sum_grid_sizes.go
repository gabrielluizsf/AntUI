package template

func sumGridSizes(sizes []int) int {
	total := 0
	for _, size := range sizes {
		total += size
	}
	return total
}
