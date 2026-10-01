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
		p.CubicTo(c1.X, c1.Y, c2.X, c2.Y, end.X, end.Y)
		theta = next
	}
}

// arcCenter turns the arc's endpoint description into the one a curve can be
// drawn from: the middle of the ellipse, the point on it the arc starts at
// measured in the ellipse's own frame, and the angle it sweeps. The middle
// sits on the perpendicular of the chord, and the large flag picks the side:
// the same side gives the arc bowing less than half a turn, the other more.
func arcCenter(p0, p1 Point, rx, ry, rot float64, large, sweep bool) (center, from Point, sweepAngle float64) {
	phi := rot * math.Pi / 180
	cos, sin := math.Cos(phi), math.Sin(phi)
	dx, dy := (p0.X-p1.X)/2, (p0.Y-p1.Y)/2
	ex, ey := cos*dx+sin*dy, -sin*dx+cos*dy

	if over := (ex*ex)/(rx*rx) + (ey*ey)/(ry*ry); over > 1 {
		rx, ry = rx*math.Sqrt(over), ry*math.Sqrt(over)
	}
	sign := 1.0
	if large == sweep {
		sign = -1
	}
	num := math.Max(rx*rx*ry*ry-rx*rx*ey*ey-ry*ry*ex*ex, 0)
	den := math.Max(rx*rx*ey*ey+ry*ry*ex*ex, 1e-12)
	coef := sign * math.Sqrt(num/den)
	center = Point{
		X: cos*coef*rx*ey/ry - sin*coef*ry*ex/rx + (p0.X+p1.X)/2,
		Y: sin*coef*rx*ey/ry + cos*coef*ry*ex/rx + (p0.Y+p1.Y)/2,
	}
	from = Point{X: (ex - coef*rx*ey/ry) / rx, Y: (ey + coef*ry*ex/rx) / ry}
	to := Point{X: (-ex - coef*rx*ey/ry) / rx, Y: (-ey + coef*ry*ex/rx) / ry}
	sweepAngle = math.Mod(turnAngle(from, to), 2*math.Pi)
	if !sweep && sweepAngle > 0 {
		sweepAngle -= 2 * math.Pi
	}
	if sweep && sweepAngle < 0 {
		sweepAngle += 2 * math.Pi
	}
	return center, from, sweepAngle
}

// arcPoint is the point of the ellipse at the angle t, measured in its own
// frame and then turned by phi.
func arcPoint(center Point, rx, ry, phi, t float64) Point {
	cos, sin := math.Cos(phi), math.Sin(phi)
	x, y := rx*math.Cos(t), ry*math.Sin(t)
	return Point{X: center.X + cos*x - sin*y, Y: center.Y + sin*x + cos*y}
}

// arcControl is where the control point of a cubic stands in for the ellipse:
// the point itself, moved along the tangent there by alpha. The tangent comes
// out of the same rotation the point went through.
func arcControl(center Point, rx, ry, phi, t, alpha float64) Point {
	cos, sin := math.Cos(phi), math.Sin(phi)
	dx, dy := -rx*math.Sin(t), ry*math.Cos(t)
	at := arcPoint(center, rx, ry, phi, t)
	return Point{X: at.X + alpha*(cos*dx-sin*dy), Y: at.Y + alpha*(sin*dx+cos*dy)}
}

// turnAngle is the angle that turns one vector onto another, in (-π, π].
func turnAngle(from, to Point) float64 {
	return math.Atan2(from.X*to.Y-from.Y*to.X, from.X*to.X+from.Y*to.Y)
}
