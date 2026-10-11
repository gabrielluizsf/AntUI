package svg

import "math"

// angle is the direction `orient="auto"` turns a marker to at this vertex: the
// way the path is going where there is only one segment beside it, and the
// middle of the turn where there are two. The middle is the two directions
// added as unit vectors — the bisector of the angle they make — which is what
// keeps a marker on a corner pointing along the way the path is turning rather
// than down one of the two sides of it. A path that turns straight back on
// itself has no middle to stand in, and the direction it came in from is the
// answer then.
func (v markerVertex) angle() float64 {
	switch {
	case v.hasIn && v.hasOut:
		ix, iy := unit(v.inX, v.inY)
		ox, oy := unit(v.outX, v.outY)
		sx, sy := ix+ox, iy+oy
		if sx == 0 && sy == 0 {
			return math.Atan2(iy, ix)
		}
		return math.Atan2(sy, sx)
	case v.hasOut:
		return math.Atan2(v.outY, v.outX)
	default:
		return math.Atan2(v.inY, v.inX)
	}
}
