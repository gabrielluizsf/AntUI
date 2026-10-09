package css

import (
	"fmt"
	"github.com/gabrielluizsf/antui/canvas"
)

// parseOKLCH reads the polar oklch() form of oklab.
func parseOKLCH(s string) (canvas.Color, error) {
	slots, alpha, hasAlpha, ok := splitColorArgs(s)
	if !ok {
		return 0, fmt.Errorf("css: %q is not an oklch() colour", s)
	}
	if len(slots) != 3 {
		return 0, fmt.Errorf("css: %q has the wrong channel count", s)
	}
	l, err := parseOKLabLight(slots[0])
	if err != nil {
		return 0, err
	}
	c, err := parseOKLCHChroma(slots[1])
	if err != nil {
		return 0, err
	}
	h, err := parseAngle(slots[2])
	if err != nil {
		return 0, err
	}
	a, b := polarToLab(c, h)
	r, g, b := okLabToSRGB(l, a, b)
	return finishOK(r, g, b, alpha, hasAlpha)
}
