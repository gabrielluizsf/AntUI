package css

import (
	"math"
	"strconv"
	"strings"
)

// parseIterations reads one iteration count: a number or the infinite
// keyword, which lands as positive infinity.
func parseIterations(raw string) (float64, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "infinite" {
		return math.Inf(1), true
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, false
	}
	return v, true
}
