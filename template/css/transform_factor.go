package css

import (
	"strings"
)

// transformFactor is a scale argument: a unitless number or a percentage,
// which divides by a hundred the way scale(150%) means scale(1.5).
func transformFactor(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		v, ok := transformNumber(strings.TrimSuffix(s, "%"))
		return v / 100, ok
	}
	return transformNumber(s)
}
