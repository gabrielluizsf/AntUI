package css

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// parseHWB reads hwb(hue white blackness / alpha): the hue angle, then two
// percentages that mix white and black into the hue.
func parseHWB(s string) (canvas.Color, error) {
	slots, alpha, hasAlpha, ok := splitColorArgs(s)
	if !ok {
		return 0, fmt.Errorf("css: %q is not an hwb() colour", s)
	}
	if len(slots) != 3 {
		return 0, fmt.Errorf("css: %q has the wrong channel count", s)
	}
	h, err := parseAngle(slots[0])
	if err != nil {
		return 0, err
	}
	w, err := parsePercent(slots[1])
	if err != nil {
		return 0, err
	}
	b, err := parsePercent(slots[2])
	if err != nil {
		return 0, err
	}
	a := uint8(255)
	if hasAlpha {
		if a, err = parseColorAlpha(alpha); err != nil {
			return 0, err
		}
	}
	c := hwbToRGB(h, w, b)
	return canvas.RGBA(c[0], c[1], c[2], a), nil
}
