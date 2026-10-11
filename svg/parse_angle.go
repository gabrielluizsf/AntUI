package svg

import (
	"math"
	"strconv"
	"strings"
)

// parseAngle reads an angle in degrees, written either as a plain number or
// with any of the units a drawing may use, which is what a hue, a rotation and
// a shear all are.
func parseAngle(s string) (float64, error) {
	s = strings.TrimSpace(s)
	// The longest unit is tried first, because one of them ends in another:
	// `grad` is `rad` with a `g` in front of it, and taking the short one would
	// leave a number with a letter stuck to the end of it.
	for _, unit := range []struct {
		suffix string
		scale  float64
	}{
		{"grad", 360.0 / 400},
		{"turn", 360},
		{"rad", 180 / math.Pi},
		{"deg", 1},
	} {
		body, ok := strings.CutSuffix(s, unit.suffix)
		if !ok {
			continue
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(body), 64)
		if err != nil {
			// It ends in something that only looks like this unit, so try the
			// next one rather than giving up on the number.
			continue
		}
		return v * unit.scale, nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}
