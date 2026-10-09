package css

import (
	"strconv"
	"strings"
)

// originAxisLength is the shared length/percentage tail of an origin axis.
func originAxisLength(s string, ctx Units) (Length, bool) {
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil {
			return Length{}, false
		}
		return Pct(v), true
	}
	if s == "left" || s == "right" || s == "top" || s == "bottom" {
		return Length{}, false
	}
	l, err := parseLengthAt(s, ctx)
	if err != nil || l.Auto() || l.None() {
		return Length{}, false
	}
	return l, true
}
