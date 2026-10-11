package canvas

import "math"

// conicGradientT walks one revolution around the centre, 0 at the start angle
// (which points up by default) and 1 back around clockwise.
func conicGradientT(x, y float64, w, h int, g Gradient) float64 {
	cx := g.CenterX * float64(w)
	cy := g.CenterY * float64(h)
	ang := math.Atan2(x-cx, -(y-cy)) - g.Angle
	ang = math.Mod(ang, 2*math.Pi)
	if ang < 0 {
		ang += 2 * math.Pi
	}
	return ang / (2 * math.Pi)
}
