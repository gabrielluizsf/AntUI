package canvas

import "math"

// turnAngle is the angle that turns one vector onto another, in (-π, π].
func turnAngle(from, to Point) float64 {
	return math.Atan2(from.X*to.Y-from.Y*to.X, from.X*to.X+from.Y*to.Y)
}
