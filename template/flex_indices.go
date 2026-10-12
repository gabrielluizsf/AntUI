package template

// flexIndices lists the tracks that carry a flex factor.
func flexIndices(plans []gridTrackPlan) []int {
	out := make([]int, 0, len(plans))
	for i := range plans {
		if plans[i].flex > 0 {
			out = append(out, i)
		}
	}
	return out
}
