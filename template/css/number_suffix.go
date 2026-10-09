package css

import (
	"math"
	"strconv"
	"strings"
)

// numberSuffix parses the float before a unit suffix. It requires the whole
// string to be number+suffix, so a "turn" does not leak into "deg" and "s"
// does not leak out of "ms" (ms is checked first by callers that care).
func numberSuffix(s, suffix string) (float64, bool) {
	if !strings.HasSuffix(s, suffix) {
		return 0, false
	}
	body := strings.TrimSpace(strings.TrimSuffix(s, suffix))
	if body == "" {
		return 0, false
	}
	v, err := strconv.ParseFloat(body, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, false
	}
	return v, true
}
