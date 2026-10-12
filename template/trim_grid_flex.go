package template

// trimGridFlex takes one unit at a time back from the tracks that can still
// shrink, in the order the caller sorted them, until the overshoot is gone,
// and returns what is left.
func trimGridFlex(plans []gridTrackPlan, indices []int, remaining int) int {
	for remaining < 0 {
		changed := false
		for _, index := range indices {
			if plans[index].size <= plans[index].min {
				continue
			}
			plans[index].size--
			remaining++
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
