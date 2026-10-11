package canvas

// Empty reports whether the path holds no subpath, or only subpaths that never
// got a point. A path with nothing in it draws nothing at all, and saying so
// early saves walking it.
func (p *Path) Empty() bool {
	for _, s := range p.subs {
		if len(s.pts) > 0 {
			return false
		}
	}
	return true
}
