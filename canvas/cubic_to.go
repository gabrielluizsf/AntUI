package canvas

// CubicTo adds a cubic Bézier using grouped control and end points.
func (p *Path) CubicTo(c1, c2, end Point) {
	p.curveTo(cubicPoints(p.end(), c1, c2, end))
}
