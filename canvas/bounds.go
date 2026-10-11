package canvas

// Bounds is the tightest box holding every point, and reports whether the path
// has any point at all to bound.
func (p *Path) Bounds() (minX, minY, maxX, maxY float64, ok bool) {
	for _, s := range p.subs {
		for _, pt := range s.pts {
			if !ok {
				minX, minY, maxX, maxY, ok = pt.X, pt.Y, pt.X, pt.Y, true
				continue
			}
			minX = min(minX, pt.X)
			minY = min(minY, pt.Y)
			maxX = max(maxX, pt.X)
			maxY = max(maxY, pt.Y)
		}
	}
	return minX, minY, maxX, maxY, ok
}
