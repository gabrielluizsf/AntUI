package canvas

import "math"

// reach is how far a control polygon reaches, which is the size of the curve it
// describes.
func reach(points ...Point) float64 {
	d := 0.0
	for i := 1; i < len(points); i++ {
		d = max(d, math.Hypot(points[i].X-points[i-1].X, points[i].Y-points[i-1].Y))
	}
	return d
}
