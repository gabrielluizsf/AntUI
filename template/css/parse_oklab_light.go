package css

import (
	"strings"
)

// parseOKLabLight reads oklab/oklch lightness: 0-1, with 100% meaning 1.
func parseOKLabLight(f string) (float64, error) {
	if strings.HasSuffix(strings.TrimSpace(f), "%") {
		v, err := parsePercent(f)
		return v / 100, err
	}
	return parsePercent(f)
}
