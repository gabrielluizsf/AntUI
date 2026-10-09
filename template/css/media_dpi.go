package css

import (
	"strconv"
	"strings"
)

// mediaDpi reads a resolution — 96dpi, 40dpcm, 2dppx — and answers it in dots
// per inch, the unit every query is compared in. A bare number, a value with
// no digits in it, or a unit the spec does not give a resolution returns
// ok=false and the caller warns.
func mediaDpi(s string) (float64, bool) {
	t := strings.TrimSpace(strings.ToLower(s))
	i := 0
	for i < len(t) && (t[i] == '.' || (t[i] >= '0' && t[i] <= '9')) {
		i++
	}
	v, err := strconv.ParseFloat(t[:i], 64)
	if err != nil {
		return 0, false
	}
	switch strings.TrimSpace(t[i:]) {
	case "dpi":
	case "dpcm":
		v *= 2.54
	case "dppx":
		v *= 96
	default:
		return 0, false
	}
	return roundDpi(v), true
}
