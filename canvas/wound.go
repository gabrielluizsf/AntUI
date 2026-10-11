package canvas

import "slices"

// wound returns a piece ordered so that it comes out the same way round as the
// rectangles it overlaps, since a piece that came out the other way would take
// its overlap away instead of adding to it.
func wound(pts []Point) []Point {
	if len(pts) > 2 && turn(pts) > 0 {
		slices.Reverse(pts)
	}
	return pts
}
