package css

import (
	"strconv"
)

// parseFlexNumber reads a non-negative flex factor such as the grow or
// shrink multiplier; a negative number is invalid, as in CSS.
func parseFlexNumber(raw string) (float64, bool) {
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil || v < 0 {
		return 0, false
	}
	return v, true
}
