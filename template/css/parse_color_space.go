package css

import (
	"fmt"
	"github.com/gabrielluizsf/antui/canvas"
	"strings"
)

// parseColorSpace reads color(space r g b / alpha), mapping the predefined
// spaces this engine can honestly paint: srgb, srgb-linear, display-p3, xyz
// and xyz-d65, plus xyz-d50 through a Bradford adaptation.
func parseColorSpace(s string) (canvas.Color, error) {
	slots, alpha, hasAlpha, ok := splitColorArgs(s)
	if !ok {
		return 0, fmt.Errorf("css: %q is not a color() colour", s)
	}
	if len(slots) != 4 {
		return 0, fmt.Errorf("css: %q has the wrong channel count", s)
	}
	space := strings.ToLower(strings.TrimSpace(slots[0]))
	r, g, b, err := parseColorComponents(slots[1:])
	if err != nil {
		return 0, err
	}
	var lr, lg, lb float64
	switch space {
	case "srgb":
		return finishOK(r, g, b, alpha, hasAlpha)
	case "srgb-linear":
		lr, lg, lb = r, g, b
	default:
		var x, y, z float64
		switch space {
		case "display-p3":
			lr, lg, lb = displayP3ToSRGB(r, g, b)
		case "xyz", "xyz-d65":
			x, y, z = r, g, b
			lr, lg, lb = xyzToLinearSRGB(x, y, z)
		case "xyz-d50":
			lr, lg, lb = xyzD50ToSRGB(r, g, b)
		default:
			return 0, fmt.Errorf("css: unsupported colour space %q", space)
		}
	}
	return finishOK(linearToSRGB(lr), linearToSRGB(lg), linearToSRGB(lb), alpha, hasAlpha)
}
