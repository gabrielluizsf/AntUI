package template

func sumGridPlans(plans []gridTrackPlan) int {
	total := 0
	for _, plan := range plans {
		total += plan.size
	}
	return total
}
