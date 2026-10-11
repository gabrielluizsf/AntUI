package canvas

import "math"

// wrap is a distance along a pattern brought back inside it, counting from the
// start even when the distance is given as a negative one: an offset of minus
// half a period is where half a period before the end is, not before the start.
func wrap(d, period float64) float64 {
	if d = math.Mod(d, period); d < 0 {
		d += period
	}
	return d
}
