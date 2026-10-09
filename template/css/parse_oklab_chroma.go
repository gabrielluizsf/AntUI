package css

import (
	"strings"
)

// parseOKLabChroma reads oklab a/b: a number, or a percentage where 100% is
// 0.4 (the OKLab reference extent).
func parseOKLabChroma(f string) (float64, error) {
	if strings.HasSuffix(strings.TrimSpace(f), "%") {
		v, err := parsePercent(f)
		return v * 0.4 / 100, err
	}
	return parsePercent(f)
}
