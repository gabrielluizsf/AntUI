package canvas

// AddPolyline adds a run of points as one subpath, joined by straight lines.
// close joins the last point back to the first, which is the difference
// between a polyline and a polygon.
func (p *Path) AddPolyline(points []Point, close bool) {
	if len(points) == 0 {
		return
	}
	p.MoveTo(points[0].X, points[0].Y)
	for _, pt := range points[1:] {
		p.LineTo(pt.X, pt.Y)
	}
	if close {
		p.Close()
	}
}
