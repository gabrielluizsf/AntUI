package font

import (
	"math"
)

// LineHeight is ascent + descent + gap, in pixels, rounded away from zero so
// that two lines never overlap by a fraction.
func (f *TTF) LineHeight(pixels float64) int {
	s := f.Scaled(pixels)
	h := (float64(f.Ascent) - float64(f.Descent) + float64(f.LineGap)) * s
	return int(math.Ceil(h))
}
