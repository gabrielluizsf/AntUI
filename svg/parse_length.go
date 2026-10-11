package svg

import (
	"strconv"
	"strings"
)

// parseLength reads a length, which is a number with an optional unit. The units
// a drawing may use are the pixels and the unitless numbers of a viewBox, and
// the absolute ones, which are turned into pixels at ninety-six to the inch, as
// CSS measures them.
func parseLength(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	for _, unit := range lengthUnits {
		if body, ok := strings.CutSuffix(s, unit.suffix); ok {
			v, err := strconv.ParseFloat(body, 64)
			if err != nil {
				// It ends in something that only looks like this unit, so the
				// next one is worth trying before giving up on the number.
				continue
			}
			return v * unit.scale, true
		}
	}
	// No unit at all is a number of the drawing's own units, which is what a
	// viewBox is written in.
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}
