package svg

import (
	"strconv"
	"strings"
)

// colorChannel reads one colour channel that is either 0-255 or a percentage of
// 255.
func colorChannel(s string) (float64, bool) {
	if p, ok := strings.CutSuffix(s, "%"); ok {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return 0, false
		}
		return v * 255 / 100, true
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

// colorAlpha reads the alpha, which is a fraction of one or a percentage of it,
// and rounds it to a byte.
func colorAlpha(s string) (uint8, bool) {
	if p, ok := strings.CutSuffix(s, "%"); ok {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return 0, false
		}
		return byte(clampFloat(v*255/100, 0, 255) + 0.5), true
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, false
	}
	return byte(clampFloat(v*255, 0, 255) + 0.5), true
}
