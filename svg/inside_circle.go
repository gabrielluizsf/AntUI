package svg

import "math"

// insideCircle pulls a focus back inside the circle it belongs to. A focus outside
// it has no highlight to give, and SVG puts it on the edge rather than refusing
// to paint the shape at all.
func insideCircle(f, c point, r float64) point {
	d := math.Hypot(f.X-c.X, f.Y-c.Y)
	if d <= r || d == 0 {
		return f
	}
	return point{X: c.X + (f.X-c.X)*r/d, Y: c.Y + (f.Y-c.Y)*r/d}
}
