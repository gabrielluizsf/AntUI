package svg

import (
	"strconv"
	"strings"
)

// parseAlpha reads an opacity, which is a number from zero to one, a percentage,
// or nothing at all, in which case it is one.
func parseAlpha(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	if p, ok := strings.CutSuffix(s, "%"); ok {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return 0, false
		}
		return clampFloat(v/100, 0, 1), true
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return clampFloat(v, 0, 1), true
}
