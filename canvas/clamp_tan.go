package canvas

import "math"

// clampTan keeps a shear angle just off the vertical.
func clampTan(rad float64) float64 {
	const limit = math.Pi/2 - 1e-3
	if rad > limit {
		return limit
	}
	if rad < -limit {
		return -limit
	}
	return rad
}
