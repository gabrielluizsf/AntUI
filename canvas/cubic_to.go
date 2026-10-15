package canvas

// CubicTo adds a cubic Bézier using grouped control and end points.
func (p *Path) CubicTo(c1, c2, end Point) {
	p.curveTo(cubicPoints(p.end(), Point{X: c1.X, Y: c1.Y}, Point{X: c2.X, Y: c2.Y}, Point{X: end.X, Y: end.Y}))
}