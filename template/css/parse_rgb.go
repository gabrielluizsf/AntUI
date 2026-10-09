package css

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// parseRGB reads rgb()/rgba() in both the legacy comma form and the modern
// space-separated form with an optional "/ alpha" divider. Channels are
// integers 0-255 or percentages; the alpha is a number 0-1 or a percentage.
func parseRGB(s string) (canvas.Color, error) {
	slots, alpha, hasAlpha, ok := splitColorArgs(s)
	if !ok {
		return 0, fmt.Errorf("css: %q is not an rgb() colour", s)
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
	r, err := rgbChannel(slots[0])
	if err != nil {
		return 0, err
	}
	g, err := rgbChannel(slots[1])
	if err != nil {
		return 0, err
	}
	b, err := rgbChannel(slots[2])
	if err != nil {
		return 0, err
	}
	return canvas.RGBA(r, g, b, a), nil
}
