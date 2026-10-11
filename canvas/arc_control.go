package canvas

import "math"

// arcControl is where the control point of a cubic stands in for the ellipse:
// the point itself, moved along the tangent there by alpha. The tangent comes
// out of the same rotation the point went through.
func arcControl(center Point, rx, ry, phi, t, alpha float64) Point {
	cos, sin := math.Cos(phi), math.Sin(phi)
	dx, dy := -rx*math.Sin(t), ry*math.Cos(t)
	at := arcPoint(center, rx, ry, phi, t)
	return Point{X: at.X + alpha*(cos*dx-sin*dy), Y: at.Y + alpha*(sin*dx+cos*dy)}
}
