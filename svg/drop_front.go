package svg

// dropFront takes the first k values off a list of places, which is what goes
// with the k letters dropped from the front of the writing.
func dropFront(list []float64, k int) []float64 {
	if k <= 0 {
		return list
	}
	if len(list) > k {
		return list[k:]
	}
	return nil
}
