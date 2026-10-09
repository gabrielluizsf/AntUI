package css

import (
	"strings"
)

// parseLCHChroma reads lch chroma: a number, or a percentage where 100% is
// 150 (the CIE reference saturation).
func parseLCHChroma(f string) (float64, error) {
	if strings.HasSuffix(strings.TrimSpace(f), "%") {
		v, err := parsePercent(f)
		return v * 1.5, err
	}
	return parsePercent(f)
}
