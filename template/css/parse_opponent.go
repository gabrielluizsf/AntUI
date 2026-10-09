package css

import (
	"strings"
)

// parseOpponent reads lab a/b: numbers, or percentages where 100% is 125.
func parseOpponent(f string) (float64, error) {
	if strings.HasSuffix(strings.TrimSpace(f), "%") {
		v, err := parsePercent(f)
		return v * 1.25, err
	}
	return parsePercent(f)
}
