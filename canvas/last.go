package canvas

// last returns the subpath being drawn, or nil when the path is empty.
func (p *Path) last() *subpath {
	if len(p.subs) == 0 {
		return nil
	}
	return &p.subs[len(p.subs)-1]
}
