package canvas

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
