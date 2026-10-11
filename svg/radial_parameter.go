package svg

import "math"

// radialParameter is where a point lands between the focus of a radial gradient,
// where its first stop is, and the circle the last stop is on. The two are not
// the same distance apart from a point that is off to one side, so the point is
// measured against where the ray from the focus through it leaves the circle
// rather than against the radius — which is what gives a focus away from the
// centre its highlight instead of a ring around the middle.
func (s sampler) radialParameter(x, y float64) float64 {
	if s.r <= 0 {
		// A gradient with no radius is a single colour, which is the last one.
		return 1
	}
	dx, dy := x-s.focus.X, y-s.focus.Y
	d := math.Hypot(dx, dy)
	if d == 0 {
		return 0
	}
	if span, ok := circleFrom(s.focus, dx/d, dy/d, s.c, s.r); ok {
		return d / span
	}
	// The ray misses the circle, which cannot happen with a focus inside it but
	// can with a radius so small the arithmetic does not hold: the distance from
	// the focus alone is as close as the point can be placed.
	return d / s.r
}
