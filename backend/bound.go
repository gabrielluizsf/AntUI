package backend

// bound holds v between low and high, where zero on either side means no
// bound there. High is applied last, so a maximum below the minimum wins.
func bound(v, low, high int) int {
	if low > 0 && v < low {
		v = low
	}
	if high > 0 && v > high {
		v = high
	}
	return v
}
