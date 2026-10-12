package template

// stretchIndices lists the tracks allowed to grow into extra space.
func stretchIndices(plans []gridTrackPlan) []int {
	out := make([]int, 0, len(plans))
	for i := range plans {
		if plans[i].stretch {
			out = append(out, i)
		}
	}
	return out
}
