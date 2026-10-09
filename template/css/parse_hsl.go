package css

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// parseHSL reads hsl()/hsla() in either syntax. The hue is an angle unit
// (deg/rad/grad/turn, or a bare number meaning degrees); saturation and
// lightness are percentages.
func parseHSL(s string) (canvas.Color, error) {
	slots, alpha, hasAlpha, ok := splitColorArgs(s)
	if !ok {
		return 0, fmt.Errorf("css: %q is not an hsl() colour", s)
	}
	a := uint8(255)
	switch {
	case hasAlpha:
		var err error
		if a, err = parseColorAlpha(alpha); err != nil {
			return 0, err
		}
	case len(slots) == 4:
		var err error
		if a, err = parseColorAlpha(slots[3]); err != nil {
			return 0, err
		}
		slots = slots[:3]
	}
	if len(slots) != 3 {
		return 0, fmt.Errorf("css: %q has the wrong channel count", s)
	}
	h, err := parseAngle(slots[0])
	if err != nil {
		return 0, err
	}
	sp, err := parsePercent(slots[1])
	if err != nil {
		return 0, err
	}
	lp, err := parsePercent(slots[2])
	if err != nil {
		return 0, err
	}
	c := hslToRGB(h, sp, lp)
	return canvas.RGBA(c[0], c[1], c[2], a), nil
}
