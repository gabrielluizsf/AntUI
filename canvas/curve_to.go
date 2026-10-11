package canvas

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
