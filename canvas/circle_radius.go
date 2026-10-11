package canvas

import "math"

// circleRadius is how far a circular ending shape of the given size reaches.
func circleRadius(cx, cy float64, w, h int, size GradientSize) float64 {
	switch size {
	case GSClosestSide:
		return min(min(cx, float64(w)-cx), min(cy, float64(h)-cy))
	case GSFarthestSide:
		return max(max(cx, float64(w)-cx), max(cy, float64(h)-cy))
	case GSClosestCorner:
		ex, ey := nearestCorner(cx, cy, w, h)
		return math.Hypot(ex-cx, ey-cy)
	default:
		ex, ey := farthestCorner(cx, cy, w, h)
		return math.Hypot(ex-cx, ey-cy)
	}
}
