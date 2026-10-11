package canvas

// QuadTo adds a quadratic Bézier from the current point to x, y, bulging
// towards the control point cx, cy. The curve is flattened on the way in, so
// the whole path stays a polygon set.
func (p *Path) QuadTo(cx, cy, x, y float64) {
	p.curveTo(quadPoints(p.end(), Point{cx, cy}, Point{x, y}))
}
