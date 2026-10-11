package canvas

// Close joins the last subpath back to its first point and opens a new subpath
// at that same point, so a line written after the close carries on from where
// the path closed — which is what a moveto-less "Z L" pair means.
func (p *Path) Close() {
	s := p.last()
	if s == nil || len(s.pts) == 0 {
		return
	}
	s.closed = true
	p.subs = append(p.subs, subpath{pts: []Point{s.pts[0]}})
}
