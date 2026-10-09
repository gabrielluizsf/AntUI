package css

import (
	"fmt"
	"github.com/gabrielluizsf/antui/canvas"
)

// parseLCH reads the polar lch() form of lab: lightness, chroma and a hue angle.
func parseLCH(s string) (canvas.Color, error) {
	slots, alpha, hasAlpha, ok := splitColorArgs(s)
	if !ok {
		return 0, fmt.Errorf("css: %q is not an lch() colour", s)
	}
	if len(slots) != 3 {
		return 0, fmt.Errorf("css: %q has the wrong channel count", s)
	}
	l, err := parseLabLight(slots[0])
	if err != nil {
		return 0, err
	}
	c, err := parseLCHChroma(slots[1])
	if err != nil {
		return 0, err
	}
	h, err := parseAngle(slots[2])
	if err != nil {
		return 0, err
	}
	a, b := polarToLab(c, h)
	x, y, z := labToXYZ(l, a, b)
	return finishColor(x, y, z, alpha, hasAlpha)
}
