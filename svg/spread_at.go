package svg

import "math"

// spreadAt is the parameter after the spread method has had its say about
// everything past the ends of the line.
func (s sampler) spreadAt(t float64) float64 {
	switch s.spread {
	case spreadRepeat:
		// Mod keeps the ends out of it, which is what makes a repeated gradient
		// start again at the first colour instead of jumping to the last one.
		return t - math.Floor(t)
	case spreadReflect:
		t = math.Mod(math.Abs(t), 2)
		if t > 1 {
			return 2 - t
		}
		return t
	}
	return t
}
