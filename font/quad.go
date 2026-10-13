package font

import "math"

// quad flattens a quadratic curve into lines.
//
// The number of pieces comes from how far the control point is from the
// straight line between the ends: a curve that is nearly straight gets one
// piece and a tight corner gets many. Fixing the count instead would either
// draw far too many at small sizes or show corners at large ones.
func (r *raster) quad(x0, y0, cx, cy, x1, y1 float64) {
	dx := (x0+x1)/2 - cx
	dy := (y0+y1)/2 - cy
	steps := int(math.Ceil(math.Sqrt(math.Hypot(dx, dy) * 2)))
	steps = min(max(steps, 1), 32)

	px, py := x0, y0
	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		u := 1 - t
		qx := u*u*x0 + 2*u*t*cx + t*t*x1
		qy := u*u*y0 + 2*u*t*cy + t*t*y1
		r.line(px, py, qx, qy)
		px, py = qx, qy
	}
}
