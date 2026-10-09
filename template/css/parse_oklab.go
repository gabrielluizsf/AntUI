package css

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// parseOKLab reads oklab(L a b / alpha) in the OKLab perceptual space.
func parseOKLab(s string) (canvas.Color, error) {
	slots, alpha, hasAlpha, ok := splitColorArgs(s)
	if !ok {
		return 0, fmt.Errorf("css: %q is not an oklab() colour", s)
	}
	if len(slots) != 3 {
		return 0, fmt.Errorf("css: %q has the wrong channel count", s)
	}
	l, err := parseOKLabLight(slots[0])
	if err != nil {
		return 0, err
	}
	a, err := parseOKLabChroma(slots[1])
	if err != nil {
		return 0, err
	}
	b, err := parseOKLabChroma(slots[2])
	if err != nil {
		return 0, err
	}
	r, g, b := okLabToSRGB(l, a, b)
	return finishOK(r, g, b, alpha, hasAlpha)
}
