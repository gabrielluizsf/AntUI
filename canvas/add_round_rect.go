package canvas

import "math"

// AddRoundRect adds a rectangle whose corners are quarter ellipses. A radius of
// zero in an axis leaves a square corner there, and a radius wider than half
// the box is pulled in to it, which is what CSS border-radius and the rounded
// corners of a shape both do.
func (p *Path) AddRoundRect(x, y, w, h, rx, ry float64) {
	rx, ry = min(math.Abs(rx), w/2), min(math.Abs(ry), h/2)
	if w <= 0 || h <= 0 {
		return
	}
	if rx <= 0 || ry <= 0 {
		p.AddRect(x, y, w, h)
		return
	}
	k := kappa
	p.MoveTo(x+rx, y)
	p.LineTo(x+w-rx, y)
	p.CubicTo(x+w-rx+k*rx, y, x+w, y+ry-k*ry, x+w, y+ry)
	p.LineTo(x+w, y+h-ry)
	p.CubicTo(x+w, y+h-ry+k*ry, x+w-rx+k*rx, y+h, x+w-rx, y+h)
	p.LineTo(x+rx, y+h)
	p.CubicTo(x+rx-k*rx, y+h, x, y+h-ry+k*ry, x, y+h-ry)
	p.LineTo(x, y+ry)
	p.CubicTo(x, y+ry-k*ry, x+rx-k*rx, y, x+rx, y)
	p.Close()
}
