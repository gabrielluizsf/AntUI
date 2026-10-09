package css

import (
	"strconv"
	"strings"
)

// scanFloat reads a floating point number with fmt-equivalent strictness.
func scanFloat(s string, dst *float64) (int, error) {
	t := strings.TrimSuffix(s, "%")
	v, err := strconv.ParseFloat(strings.TrimSpace(t), 64)
	if err != nil {
		return 0, err
	}
	*dst = v
	if strings.HasSuffix(strings.TrimSpace(s), "%") {
		*dst = v / 100
	}
	return len(t), nil
}
