package canvas

import "math"

// arcPoint is the point of the ellipse at the angle t, measured in its own
// frame and then turned by phi.
func arcPoint(center Point, rx, ry, phi, t float64) Point {
	cos, sin := math.Cos(phi), math.Sin(phi)
	x, y := rx*math.Cos(t), ry*math.Sin(t)
	return Point{X: center.X + cos*x - sin*y, Y: center.Y + sin*x + cos*y}
}
