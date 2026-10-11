package canvas

// quadPoints flattens a quadratic Bézier into straight segments.
func quadPoints(p0, c, p1 Point) []Point {
	n := curveSegments(reach(p0, c, p1))
	pts := make([]Point, 0, n)
	for i := 1; i <= n; i++ {
		t := float64(i) / float64(n)
		u := 1 - t
		pts = append(pts, Point{
			X: u*u*p0.X + 2*u*t*c.X + t*t*p1.X,
			Y: u*u*p0.Y + 2*u*t*c.Y + t*t*p1.Y,
		})
	}
	return pts
}
