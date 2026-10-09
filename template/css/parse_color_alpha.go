package css

import (
	"strconv"
	"strings"
)

// parseColorAlpha reads the "/ 0.5" or "/ 50%" part of a colour: a number
// 0-1 or a percentage, clamped and scaled to an alpha byte.
func parseColorAlpha(f string) (uint8, error) {
	f = strings.TrimSpace(f)
	if strings.HasSuffix(f, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(f, "%"), 64)
		if err != nil {
			return 0, err
		}
		return uint8(clamp(v/100)*255 + 0.5), nil
	}
	v, err := strconv.ParseFloat(f, 64)
	if err != nil {
		return 0, err
	}
	return uint8(clamp(v)*255 + 0.5), nil
}
