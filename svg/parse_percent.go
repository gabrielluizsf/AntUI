package svg

import (
	"fmt"
	"strconv"
	"strings"
)

// parsePercent reads a number that must be written with its sign, as the
// saturation and lightness of a colour are.
func parsePercent(s string) (float64, error) {
	body, ok := strings.CutSuffix(strings.TrimSpace(s), "%")
	if !ok {
		return 0, fmt.Errorf("antui/svg: %q is not a percentage", s)
	}
	return strconv.ParseFloat(strings.TrimSpace(body), 64)
}
