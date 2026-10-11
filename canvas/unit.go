package canvas

import "math"

// unit is a vector of the given length scaled to one, and the zero vector when
// it has no length at all.
func unit(x, y float64) Point {
	l := math.Hypot(x, y)
	if l == 0 {
		return Point{}
	}
	return Point{x / l, y / l}
}
