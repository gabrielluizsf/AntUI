package template

// topUpGridFlex hands one unit at a time to the tracks that can still grow,
// in the order the caller sorted them, until the pool runs out, and returns
// what is left.
func topUpGridFlex(plans []gridTrackPlan, indices []int, remaining int) int {
	for remaining > 0 {
		changed := false
		for _, index := range indices {
			if plans[index].size >= plans[index].max {
				continue
			}
			plans[index].size++
			remaining--
			changed = true
			if remaining == 0 {
				break
			}
		}
		if !changed {
			break
		}
	}
	return remaining
}
