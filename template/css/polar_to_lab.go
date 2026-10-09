package css

import (
	"math"
)

// polarToLab converts polar lch coordinates into lab a/b.
func polarToLab(c, h float64) (a, b float64) {
	h = h * math.Pi / 180
	return c * math.Cos(h), c * math.Sin(h)
}
