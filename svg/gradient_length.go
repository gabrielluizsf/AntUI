package svg

import (
	"strconv"
	"strings"
)

// gradientLength reads one number of a gradient's geometry, which is a number
// with an optional unit or a percentage. A percentage is a fraction of whatever
// the units say the numbers are measured against, which is the shape's own box —
// a viewport of one — or, in user units, the side of the drawing the number is
// on. The axis is which of those it is, and -1 is the one axis a percentage has
// no meaning for.
func gradientLength(raw string, viewport float64) (float64, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, false
	}
	if p, ok := strings.CutSuffix(raw, "%"); ok {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return 0, false
		}
		return v * viewport / 100, true
	}
	return parseLength(raw)
}
