package css

import (
	"strconv"
	"strings"
)

// parsePercent reads a percentage or a bare number, both scaled in
// percent-units.
func parsePercent(f string) (float64, error) {
	f = strings.TrimSpace(f)
	return strconv.ParseFloat(strings.TrimSuffix(f, "%"), 64)
}
