package canvas

import "math"

// dashStart is where the walk of a run begins once the offset's leftover gone
// has been spent: the segment it lands in, the point it lands on, and the part
// of that segment the walk is already past. It reports false when the run is
// spent before the walk reaches its last segment.
func dashStart(sub []Point, gone float64) (int, Point, float64, bool) {
	from := 0
	for from+1 < len(sub) {
		seg := math.Hypot(sub[from+1].X-sub[from].X, sub[from+1].Y-sub[from].Y)
		if gone < seg {
			break
		}
		gone, from = gone-seg, from+1
	}
	if from+1 >= len(sub) {
		return 0, Point{}, 0, false
	}
	at := sub[from]
	if gone > 0 {
		a, b := sub[from], sub[from+1]
		seg := math.Hypot(b.X-a.X, b.Y-a.Y)
		at = Point{a.X + (b.X-a.X)*gone/seg, a.Y + (b.Y-a.Y)*gone/seg}
	}
	return from, at, gone, true
}
