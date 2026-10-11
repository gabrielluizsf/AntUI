package svg

import (
	"strconv"
	"strings"
)

// filterAmount reads one function's argument: a number or a percentage, where
// an empty argument is the identity of one. Something else is not a number at
// all, and the caller says so instead of the filter quietly becoming full
// strength or none.
func filterAmount(args string) (float64, bool) {
	if args == "" {
		return 1, true
	}
	if p, ok := strings.CutSuffix(args, "%"); ok {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return 0, false
		}
		return v / 100, true
	}
	v, err := strconv.ParseFloat(args, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}
