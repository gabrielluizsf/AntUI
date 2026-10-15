package canvas

import "math"

// ArcParams holds parameters for an SVG elliptical arc.
type ArcParams struct {
	RX, RY     float64
	Rot        float64
	Large      bool
	Sweep      bool
	End        Point
}

// ArcTo adds the elliptical arc with grouped parameters.
func (p *Path) ArcTo(params ArcParams) {
	rx, ry, rot, large, sweep := params.RX, params.RY, params.Rot, params.Large, params.Sweep
	x, y := params.End.X, params.End.Y
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
		p.CubicTo(Point{X: c1.X, Y: c1.Y}, Point{X: c2.X, Y: c2.Y}, Point{X: end.X, Y: end.Y})
		theta = next
	}
}