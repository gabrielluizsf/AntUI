package template

// distributeGridExtra shares what is left over between the tracks that may grow
// into it, one share each and the rest to the first.
func distributeGridExtra(plans []gridTrackPlan, indices []int, extra int) {
	if len(indices) == 0 || extra <= 0 {
		return
	}
	each := extra / len(indices)
	rem := extra % len(indices)
	for i, index := range indices {
		plans[index].size += each
		if i < rem {
			plans[index].size++
		}
	}
}
