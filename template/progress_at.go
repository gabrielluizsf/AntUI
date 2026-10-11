package template

import "time"

// progressAt is elapsed over timeout clamped into the range a bar reads.
// A timeout the API cannot produce (zero or less) has no denominator, and
// the bar stays where it is rather than inventing a number.
func progressAt(elapsed, timeout time.Duration) float32 {
	if timeout <= 0 {
		return 0
	}
	v := float64(elapsed) / float64(timeout)
	switch {
	case v < 0:
		return 0
	case v > 1:
		return 1
	default:
		return float32(v)
	}
}
