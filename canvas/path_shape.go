package canvas

import "math"

// AddRect adds an axis-aligned rectangle as a closed subpath of its own.
func (p *Path) AddRect(x, y, w, h float64) {
	p.MoveTo(x, y)
	p.LineTo(x+w, y)
	p.LineTo(x+w, y+h)
	p.LineTo(x, y+h)
	p.Close()
}

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

// AddEllipse adds an ellipse centred on cx, cy as a closed subpath, drawn as
// the four quarter turns a circle is a special case of. A radius of zero or
// less adds nothing: there is no such shape to draw.
func (p *Path) AddEllipse(cx, cy, rx, ry float64) {
	if rx <= 0 || ry <= 0 {
		return
	}
	p.MoveTo(cx+rx, cy)
	p.CubicTo(cx+rx, cy+kappa*ry, cx+kappa*rx, cy+ry, cx, cy+ry)
	p.CubicTo(cx-kappa*rx, cy+ry, cx-rx, cy+kappa*ry, cx-rx, cy)
	p.CubicTo(cx-rx, cy-kappa*ry, cx-kappa*rx, cy-ry, cx, cy-ry)
	p.CubicTo(cx+kappa*rx, cy-ry, cx+rx, cy-kappa*ry, cx+rx, cy)
	p.Close()
}

// AddPolyline adds a run of points as one subpath, joined by straight lines.
// close joins the last point back to the first, which is the difference
// between a polyline and a polygon.
func (p *Path) AddPolyline(points []Point, close bool) {
	if len(points) == 0 {
		return
	}
	p.MoveTo(points[0].X, points[0].Y)
	for _, pt := range points[1:] {
		p.LineTo(pt.X, pt.Y)
	}
	if close {
		p.Close()
	}
}

// AddCircle adds a disc of the given radius, the ellipse every icon's dot is.
func (p *Path) AddCircle(cx, cy, r float64) { p.AddEllipse(cx, cy, r, r) }

// kappa is the distance a cubic's control point stands off the corner of the
// quarter turn it stands in for: the tangent of a quarter of a circle, minus
// the circle's own radius.
const kappa = 0.5522847498307936
