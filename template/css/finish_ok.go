package css

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// finishOK turns linear sRGB components into a canvas colour, applying the
// transfer curve, clamping and an optional alpha.
func finishOK(r, g, b float64, alpha string, hasAlpha bool) (canvas.Color, error) {
	var a uint8 = 255
	if hasAlpha {
		var err error
		if a, err = parseColorAlpha(alpha); err != nil {
			return 0, err
		}
	}
	return canvas.RGBA(toByte(r), toByte(g), toByte(b), a), nil
}
