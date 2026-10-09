package css

import (
	"strconv"
	"strings"
)

// rgbChannel reads one rgb() channel: an integer 0-255, or a percentage
// scaled into that range. Either is clamped.
func rgbChannel(f string) (uint8, error) {
	f = strings.TrimSpace(f)
	if strings.HasSuffix(f, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(f, "%"), 64)
		if err != nil {
			return 0, err
		}
		return clampByte(v * 255 / 100), nil
	}
	v, err := strconv.ParseFloat(f, 64)
	if err != nil {
		return 0, err
	}
	return clampByte(v), nil
}
