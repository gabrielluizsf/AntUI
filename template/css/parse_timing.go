package css

import (
	"math"
	"strconv"
	"strings"
)

// parseTiming reads one easing keyword or a cubic-bezier(x1 y1 x2 y2)
// function. The horizontal controls must stay inside [0,1] for the curve to
// be monotonic, which the parser checks the way a browser does.
func parseTiming(raw string) (Timing, bool) {
	s := strings.TrimSpace(strings.ToLower(raw))
	switch s {
	case "linear":
		return Linear, true
	case "ease":
		return Timing{x1: 0.25, y1: 0.1, x2: 0.25, y2: 1}, true
	case "ease-in":
		return Timing{x1: 0.42, y1: 0, x2: 1, y2: 1}, true
	case "ease-out":
		return Timing{x1: 0, y1: 0, x2: 0.58, y2: 1}, true
	case "ease-in-out":
		return Timing{x1: 0.42, y1: 0, x2: 0.58, y2: 1}, true
	}
	if strings.HasPrefix(s, "cubic-bezier(") && strings.HasSuffix(s, ")") {
		parts := splitFields(s[len("cubic-bezier("):len(s)-1], ',')
		if len(parts) != 4 {
			return Timing{}, false
		}
		vals := [4]float64{}
		for i, p := range parts {
			v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
			if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
				return Timing{}, false
			}
			vals[i] = v
		}
		if vals[0] < 0 || vals[0] > 1 || vals[2] < 0 || vals[2] > 1 {
			return Timing{}, false
		}
		return Timing{x1: vals[0], y1: vals[1], x2: vals[2], y2: vals[3]}, true
	}
	return Timing{}, false
}
