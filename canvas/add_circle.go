package canvas

// AddCircle adds a disc of the given radius, the ellipse every icon's dot is.
func (p *Path) AddCircle(cx, cy, r float64) { p.AddEllipse(Point{X: cx, Y: cy}, Point{X: r, Y: r}) }
