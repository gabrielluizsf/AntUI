package svg

import (
	"strconv"
	"strings"
)

// maskRegionValue is one number of a mask's region: a percentage, which is a
// fraction of the box of the element when the region is written in fractions
// of it and a fraction of the drawing when it is not; or a number, which in
// the first case is a fraction on its own and in the second is a length of the
// drawing's own. A length with units is not a fraction of anything, so it is
// only read where the region is written in the drawing's coordinates.
func maskRegionValue(raw string, boxUnits bool, span float64) (float64, bool) {
	s := strings.TrimSpace(raw)
	if p, ok := strings.CutSuffix(s, "%"); ok {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return 0, false
		}
		if boxUnits {
			return v / 100, true
		}
		return v / 100 * span, true
	}
	if boxUnits {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, false
		}
		return v, true
	}
	return parseLength(s)
}
