package canvas

// Transform maps every point of the path through m, in place. Flattening
// happens before it, so a curve rotated into place is rotated as the polygon
// it already became — which is the flattening this rasterizer is honest about.
func (p *Path) Transform(m Matrix) {
	for _, s := range p.subs {
		for i := range s.pts {
			s.pts[i].X, s.pts[i].Y = m.Map(s.pts[i].X, s.pts[i].Y)
		}
	}
}
