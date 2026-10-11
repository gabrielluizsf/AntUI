package canvas

import "math"

// arcPoints walks the circle of radius r around c from one angle to the other,
// turning the way step says, which for a join is the way the corner opens and
// for a cap is the way out from the line. A rounding is nothing but this walk,
// which is why it is a list of points and not a shape of its own.
func arcPoints(c Point, r, from, to, step float64) []Point {
	turn := to - from
	if step < 0 && turn > 0 {
		turn -= 2 * math.Pi
	}
	if step > 0 && turn < 0 {
		turn += 2 * math.Pi
	}
	// The gap between a chord and the arc it stands in for shrinks with the
	// square of the number of chords, so a long way round a big circle wants
	// proportionally more of them than a short way round a small one.
	n := max(int(math.Ceil(1.2*r*math.Abs(turn))), 2)
	pts := make([]Point, 0, n+1)
	for i := range n + 1 {
		angle := from + turn*float64(i)/float64(n)
		pts = append(pts, Point{c.X + r*math.Cos(angle), c.Y + r*math.Sin(angle)})
	}
	return pts
}
