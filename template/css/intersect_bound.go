package css

// intersectBound narrows one edge of two windows to their overlap: a minimum
// (rise) climbs to the higher of the two bounds, a maximum falls to the lower,
// and a side that declares no bound there passes the other side's through.
// It reads ranges of any ordered kind — the pixel edges are integers, the
// resolution bounds are dots per inch.
func intersectBound[T int | float64](hasA bool, a T, hasB bool, b T, rise bool) (T, bool) {
	switch {
	case hasA && hasB:
		if rise && b > a {
			return b, true
		}
		if !rise && b < a {
			return b, true
		}
		return a, true
	case hasA:
		return a, true
	case hasB:
		return b, true
	}
	var zero T
	return zero, false
}
