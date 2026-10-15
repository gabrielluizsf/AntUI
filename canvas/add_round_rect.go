package canvas

import "math"

// AddRoundRect adds a rounded rectangle as a closed subpath. The origin, size
// and corner radii are grouped into points so the signature stays compact.
func (p *Path) AddRoundRect(origin, size, radii Point) {
	x, y := origin.X, origin.Y
	w, h := size.X, size.Y
	rx, ry := radii.X, radii.Y
	rx, ry = min(math.Abs(rx), w/2), min(math.Abs(ry), h/2)
	if w <= 0 || h <= 0 {
		return
	}
	if rx <= 0 || ry <= 0 {
		p.AddRect(origin, size)
		return
	}
	k := kappa
	p.MoveTo(x+rx, y)
	p.LineTo(x+w-rx, y)
	p.CubicTo(Point{X: x + w - rx + k*rx, Y: y}, Point{X: x + w, Y: y + ry - k*ry}, Point{X: x + w, Y: y + ry})
	p.LineTo(x+w, y+h-ry)
	p.CubicTo(Point{X: x + w, Y: y + h - ry + k*ry}, Point{X: x + w - rx + k*rx, Y: y + h}, Point{X: x + w - rx, Y: y + h})
	p.LineTo(x+rx, y+h)
	p.CubicTo(Point{X: x + rx - k*rx, Y: y + h}, Point{X: x, Y: y + h - ry + k*ry}, Point{X: x, Y: y + h - ry})
	p.LineTo(x, y+ry)
	p.CubicTo(Point{X: x, Y: y + ry - k*ry}, Point{X: x + rx - k*rx, Y: y}, Point{X: x + rx, Y: y})
	p.Close()
}
