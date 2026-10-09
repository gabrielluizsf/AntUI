package css

import (
	"strconv"
	"strings"
)

// filterAmount reads a filter's argument: a percentage or a bare number,
// where the empty argument is the identity of 1.
func filterAmount(args string) float64 {
	if args == "" {
		return 1
	}
	if strings.HasSuffix(args, "%") {
		if v, err := strconv.ParseFloat(strings.TrimSuffix(args, "%"), 64); err == nil {
			return v / 100
		}
		return 1
	}
	v, err := strconv.ParseFloat(args, 64)
	if err != nil {
		return 1
	}
	return v
}
