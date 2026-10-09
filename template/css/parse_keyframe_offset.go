package css

import (
	"math"
	"strconv"
	"strings"
)

// parseKeyframeOffset reads a frame's main: from, to, or a percentage such
// as 32.5%. Percentages clamp into [0,1]; junk reports failure.
func parseKeyframeOffset(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	switch strings.ToLower(s) {
	case "from":
		return 0, true
	case "to":
		return 1, true
	}
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
			return 0, false
		}
		return clamp(v / 100), true
	}
	return 0, false
}
