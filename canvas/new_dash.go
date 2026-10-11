package canvas

import "math"

// NewDash makes a dash pattern out of the lengths stroke-dasharray was given,
// and reports whether it is a pattern at all. An odd list is doubled rather
// than refused, which is what lets a single number mean "dashes as long as the
// gaps".
func NewDash(lengths []float64, offset float64) (Dash, bool) {
	if len(lengths) == 0 {
		return Dash{}, false
	}
	pattern := append([]float64{}, lengths...)
	if len(pattern)%2 == 1 {
		pattern = append(pattern, pattern...)
	}
	// The period is what the pattern as it will be walked measures, so it is
	// counted after an odd list has been doubled: "ten on, five off" repeated
	// comes to twenty-five, not to fifteen.
	sum := 0.0
	for _, l := range pattern {
		if l <= 0 || math.IsNaN(l) {
			return Dash{}, false
		}
		sum += l
	}
	if sum <= 0 {
		return Dash{}, false
	}
	// An odd pattern is read twice over, so the dashes and the gaps are the two
	// halves of it and there is one of each of them.
	halves := len(pattern) / 2
	d := Dash{
		On:     append([]float64{}, pattern[:halves]...),
		Off:    append([]float64{}, pattern[halves:]...),
		Offset: wrap(offset, sum),
	}
	d.pattern, d.period = pattern, sum
	return d, true
}
