package svg

// dropPast keeps the values that still have a letter to sit on after the back
// of the writing has been trimmed: n letters were there, keep are now.
func dropPast(list []float64, keep, n int) []float64 {
	if keep >= n {
		return list
	}
	if len(list) > keep {
		return list[:keep]
	}
	return list
}
