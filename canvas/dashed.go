package canvas

// DashedPath is the path a dashed stroke draws: the same line cut into the runs
// that are on, each as a subpath of its own so that the two ends of a dash can
// be capped. A line that is not dashed comes back as it was, and a path of
// nothing is still a path of nothing.
func (p *Path) Dashed(d Dash) *Path {
	if p == nil {
		return nil
	}
	pattern, period := d.cycle()
	if len(pattern) == 0 || period <= 0 {
		return p
	}
	pts, closed := p.Points()
	out := NewPath()
	for i, sub := range pts {
		if len(sub) < 2 {
			if len(sub) == 1 {
				out.MoveTo(sub[0].X, sub[0].Y)
			}
			continue
		}
		if closed[i] {
			sub = append(append([]Point{}, sub...), sub[0])
		}
		// The pattern is walked from the start for every subpath, so a dash
		// that would have begun a while ago begins a while ago here too —
		// which is what a shape drawn as one path expects and what treating
		// each subpath as its own line would get wrong.
		out.appendDashes(sub, pattern, d.Offset, period)
	}
	return out
}
