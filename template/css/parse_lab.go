package css

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// parseLab reads lab(L a b / alpha) and parseLCH/parseOKLab/parseOKLCH the
// spaces alongside it. The CIE spaces are mapped to sRGB and clamped, which
// is the honest canvas behaviour for out-of-gamut colours.
func parseLab(s string) (canvas.Color, error) {
	slots, alpha, hasAlpha, ok := splitColorArgs(s)
	if !ok {
		return 0, fmt.Errorf("css: %q is not a lab() colour", s)
	}
	if len(slots) != 3 {
		return 0, fmt.Errorf("css: %q has the wrong channel count", s)
	}
	l, err := parseLabLight(slots[0])
	if err != nil {
		return 0, err
	}
	a, err := parseOpponent(slots[1])
	if err != nil {
		return 0, err
	}
	b, err := parseOpponent(slots[2])
	if err != nil {
		return 0, err
	}
	x, y, z := labToXYZ(l, a, b)
	return finishColor(x, y, z, alpha, hasAlpha)
}
