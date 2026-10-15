package canvas

// QuadTo adds a quadratic Bézier using the control and end points grouped into
// points.
func (p *Path) QuadTo(control, end Point) {
	p.curveTo(quadPoints(p.end(), control, end))
}
