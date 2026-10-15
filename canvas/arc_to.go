package canvas

import "math"

// ArcTo adds the elliptical arc of the SVG arc command: it runs from the point
// the path is at to x, y, with the two radii rx and ry, the x-axis rotated by
// rot degrees, and the two flags that choose which of the two arcs joining the
// points is meant — large for the one sweeping more than half a turn, sweep for
// the direction. Radii too small to reach the far end are grown until they do,
// which is what the arc's own definition asks for; a zero radius leaves a
// straight line, because that is the limit the arc describes.
func (p *Path) ArcTo(rx, ry, rot float64, large, sweep bool, x, y float64) {
	p0 := p.end()
	if rx == 0 || ry == 0 {
		p.LineTo(x, y)
		return
	}
	center, from, sweepAngle := arcCenter(p0, Point{x, y}, rx, ry, rot, large, sweep)
	rx, ry, phi := math.Abs(rx), math.Abs(ry), rot*math.Pi/180

	n := max(int(math.Ceil(math.Abs(sweepAngle)/(math.Pi/2))), 1)
	step := sweepAngle / float64(n)
	alpha := 4.0 / 3.0 * math.Tan(step/4)
	theta := math.Atan2(from.Y, from.X)
	for range n {
		next := theta + step
		c1 := arcControl(center, rx, ry, phi, theta, alpha)
		c2 := arcControl(center, rx, ry, phi, next, -alpha)
		end := arcPoint(center, rx, ry, phi, next)
		p.CubicTo(Point{X:c1.X,Y:c1.Y}, Point{X:c2.X,Y:c2.Y}, Point{X:end.X,Y:end.Y})
		theta = next
	}
}
