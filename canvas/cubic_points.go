package canvas

// cubicPoints flattens a cubic Bézier, the curve every arc, circle and rounded
// corner of a vector drawing is built out of.
func cubicPoints(p0, c1, c2, p1 Point) []Point {
	n := curveSegments(reach(p0, c1, c2, p1))
	pts := make([]Point, 0, n)
	for i := 1; i <= n; i++ {
		t := float64(i) / float64(n)
		u := 1 - t
		aa, bb, cc, dd := u*u*u, 3*u*u*t, 3*u*t*t, t*t*t
		pts = append(pts, Point{
			X: aa*p0.X + bb*c1.X + cc*c2.X + dd*p1.X,
			Y: aa*p0.Y + bb*c1.Y + cc*c2.Y + dd*p1.Y,
		})
	}
	return pts
}
