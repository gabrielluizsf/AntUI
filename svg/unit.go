package svg

import "math"

// unit is a direction of no length at all answered as no direction at all, so
// that a segment of zero length adds nothing to a bisector rather than adding
// a number divided by nothing.
func unit(x, y float64) (float64, float64) {
	if x == 0 && y == 0 {
		return 0, 0
	}
	n := math.Hypot(x, y)
	return x / n, y / n
}
