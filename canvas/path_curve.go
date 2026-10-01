package canvas

import "math"

// QuadTo adds a quadratic Bézier from the current point to x, y, bulging
// towards the control point cx, cy. The curve is flattened on the way in, so
// the whole path stays a polygon set.
func (p *Path) QuadTo(cx, cy, x, y float64) {
	p.curveTo(quadPoints(p.end(), Point{cx, cy}, Point{x, y}))
}

// CubicTo adds a cubic Bézier from the current point to x, y, with the two
// control points that shape it.
func (p *Path) CubicTo(c1x, c1y, c2x, c2y, x, y float64) {
	p.curveTo(cubicPoints(p.end(), Point{c1x, c1y}, Point{c2x, c2y}, Point{x, y}))
}

// curveTo appends the points of an already-flattened curve, leaving out the
// one that repeats the point the path is already at.
func (p *Path) curveTo(points []Point) {
	if len(points) == 0 {
		return
	}
	if len(p.subs) == 0 {
		p.MoveTo(points[0].X, points[0].Y)
		points = points[1:]
	}
	for _, pt := range points {
		p.LineTo(pt.X, pt.Y)
	}
}

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

// reach is how far a control polygon reaches, which is the size of the curve it
// describes.
func reach(points ...Point) float64 {
	d := 0.0
	for i := 1; i < len(points); i++ {
		d = max(d, math.Hypot(points[i].X-points[i-1].X, points[i].Y-points[i-1].Y))
	}
	return d
}

// curveSegments is how many straight segments a curve of the given size is
// flattened into. The gap between a chord and the arc it stands for shrinks
// with the square of the number of segments, so the count follows the square
// root of the curve's size: a quarter turn twenty pixels across wants about
// eight, and a tiny one wants two — anything less and a corner that should be
// a curve turns into a visible facet.
func curveSegments(d float64) int {
	return max(int(math.Ceil(1.6*math.Sqrt(max(d, 0)))), 2)
}
