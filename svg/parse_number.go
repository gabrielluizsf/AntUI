package svg

import (
	"strconv"
	"strings"
)

// parseNumber reads one plain number, which is what the path data and the point
// lists are written in.
func parseNumber(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}
