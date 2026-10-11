package template

func sumGridPlansFor(indices []int, plans []gridTrackPlan) int {
	total := 0
	for _, index := range indices {
		total += plans[index].size
	}
	return total
}
