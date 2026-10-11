package canvas

// AddRect adds an axis-aligned rectangle as a closed subpath of its own.
func (p *Path) AddRect(x, y, w, h float64) {
	p.MoveTo(x, y)
	p.LineTo(x+w, y)
	p.LineTo(x+w, y+h)
	p.LineTo(x, y+h)
	p.Close()
}
