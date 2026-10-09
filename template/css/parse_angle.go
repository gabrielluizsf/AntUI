package css

import (
	"math"
	"strconv"
	"strings"
)

// parseAngle reads a hue: a bare number means degrees; deg, rad, grad and
// turn convert to degrees.
func parseAngle(f string) (float64, error) {
	f = strings.TrimSpace(f)
	for _, c := range []struct {
		name string
		mult float64
	}{
		{"turn", 360},
		{"rad", 180 / math.Pi},
		{"grad", 0.9},
		{"deg", 1},
	} {
		if strings.HasSuffix(f, c.name) {
			v, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(f, c.name)), 64)
			if err != nil {
				return 0, err
			}
			return v * c.mult, nil
		}
	}
	return strconv.ParseFloat(f, 64)
}
