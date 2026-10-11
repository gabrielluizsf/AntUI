package canvas

// CubicTo adds a cubic Bézier from the current point to x, y, with the two
// control points that shape it.
func (p *Path) CubicTo(c1x, c1y, c2x, c2y, x, y float64) {
	p.curveTo(cubicPoints(p.end(), Point{c1x, c1y}, Point{c2x, c2y}, Point{x, y}))
}
