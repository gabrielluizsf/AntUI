package canvas

import "math"

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
