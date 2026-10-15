package canvas

// AddEllipse adds an ellipse as a closed subpath. The center and the X/Y radii
// are grouped into points so the signature stays compact.
func (p *Path) AddEllipse(center, radii Point) {
	cx, cy := center.X, center.Y
	rx, ry := radii.X, radii.Y
	if rx <= 0 || ry <= 0 {
		return
	}
	p.MoveTo(cx+rx, cy)
	p.CubicTo(Point{X:cx+rx,Y:cy+kappa*ry}, Point{X:cx+kappa*rx,Y:cy+ry}, Point{X:cx,Y:cy+ry})
	p.CubicTo(Point{X:cx-kappa*rx,Y:cy+ry}, Point{X:cx-rx,Y:cy+kappa*ry}, Point{X:cx-rx,Y:cy})
	p.CubicTo(Point{X:cx-rx,Y:cy-kappa*ry}, Point{X:cx-kappa*rx,Y:cy-ry}, Point{X:cx,Y:cy-ry})
	p.CubicTo(Point{X:cx+kappa*rx,Y:cy-ry}, Point{X:cx+rx,Y:cy-kappa*ry}, Point{X:cx+rx,Y:cy})
	p.Close()
}