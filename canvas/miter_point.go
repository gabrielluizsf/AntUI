package canvas

import "math"

// miterPoint is where the outer edge of the incoming segment meets the outer
// edge of the outgoing one, and whether they meet close enough to keep. Both
// edges run along their own segment, half a width out from it, so the corner
// between them is found by crossing those two lines — and the spike is only
// worth drawing while it stays within the miter limit, which counts the tip's
// distance from the corner in half-widths, the way SVG measures it.
func miterPoint(b, p1, p2, in, next Point, sweep, half, limit float64) (Point, bool) {
	if sweep == 0 {
		return Point{}, false
	}
	// p1 + t*in = p2 + s*next, solved for t by crossing both sides with next.
	t := ((p2.X-p1.X)*next.Y - (p2.Y-p1.Y)*next.X) / sweep
	m := Point{p1.X + in.X*t, p1.Y + in.Y*t}
	if math.Hypot(m.X-b.X, m.Y-b.Y) > limit*half {
		return Point{}, false
	}
	return m, true
}
