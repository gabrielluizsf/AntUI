package template

// clampFloat keeps v inside [0,1] for slider and progress values.
func clampFloat(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
