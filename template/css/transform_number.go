package css

import (
	"math"
	"strconv"
	"strings"
)

// transformNumber is one of matrix()'s six numbers.
func transformNumber(s string) (float64, bool) {
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, false
	}
	return v, true
}
