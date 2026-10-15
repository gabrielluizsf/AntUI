package canvas

// AddRect adds an axis-aligned rectangle as a closed subpath. The origin and
// size are grouped into points so the signature stays at two parameters.
func (p *Path) AddRect(origin, size Point) {
	x, y := origin.X, origin.Y
	w, h := size.X, size.Y
	p.MoveTo(x, y)
	p.LineTo(x+w, y)
	p.LineTo(x+w, y+h)
	p.LineTo(x, y+h)
	p.Close()
}