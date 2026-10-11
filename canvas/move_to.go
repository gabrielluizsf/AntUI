package canvas

// MoveTo starts a new subpath at x, y. A move with no line behind it still
// opens a subpath, so a path of only moves fills nothing but strokes dots.
func (p *Path) MoveTo(x, y float64) {
	p.subs = append(p.subs, subpath{pts: []Point{{x, y}}})
}
