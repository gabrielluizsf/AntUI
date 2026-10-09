package css

import (
	"strconv"
	"strings"
)

// parseColorComponents reads the three channels of color(): numbers 0-1 or
// percentages, scaled to 0-1.
func parseColorComponents(slots []string) (r, g, b float64, err error) {
	vals := [3]float64{}
	for i, f := range slots {
		f = strings.TrimSpace(f)
		if strings.HasSuffix(f, "%") {
			v, perr := strconv.ParseFloat(strings.TrimSuffix(f, "%"), 64)
			if perr != nil {
				return 0, 0, 0, perr
			}
			vals[i] = v / 100
			continue
		}
		v, perr := strconv.ParseFloat(f, 64)
		if perr != nil {
			return 0, 0, 0, perr
		}
		vals[i] = v
	}
	return vals[0], vals[1], vals[2], nil
}
