package css

import (
	"math"
)

// linearToSRGB runs the sRGB transfer curve on one linear channel.
func linearToSRGB(v float64) float64 {
	if v <= 0.0031308 {
		return 12.92 * v
	}
	return 1.055*math.Pow(v, 1.0/2.4) - 0.055
}
