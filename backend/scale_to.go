package backend

import "math"

// scaleTo multiplies a length by the display scale, for a size given in
// points. A zero or absurd scale is left alone rather than believed.
func scaleTo(v int, scale float64) int {
	if v <= 0 || scale <= 0 || scale == 1 {
		return v
	}
	return max(int(math.Round(float64(v)*scale)), 1)
}
