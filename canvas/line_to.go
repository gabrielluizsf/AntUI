package canvas

// LineTo adds a point to the subpath being drawn, starting one at the point
// itself when the path is empty — the first command of a path needs no move to
// have come before it.
func (p *Path) LineTo(x, y float64) {
	if len(p.subs) == 0 {
		p.MoveTo(x, y)
		return
	}
	s := &p.subs[len(p.subs)-1]
	s.pts = append(s.pts, Point{x, y})
}
