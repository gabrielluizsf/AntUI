package canvas

// end is the point the path is currently at.
func (p *Path) end() Point {
	s := p.last()
	if s == nil || len(s.pts) == 0 {
		return Point{}
	}
	return s.pts[len(s.pts)-1]
}
