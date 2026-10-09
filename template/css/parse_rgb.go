package css

import (
	"fmt"
	"github.com/gabrielluizsf/antui/canvas"
	"strconv"
	"strings"
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
	ch := func(f string) (uint8, error) {
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
	r, err := ch(slots[0])
	if err != nil {
		return 0, err
	}
	g, err := ch(slots[1])
	if err != nil {
		return 0, err
	}
	b, err := ch(slots[2])
	if err != nil {
		return 0, err
	}
	return canvas.RGBA(r, g, b, a), nil
}
