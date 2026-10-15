package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// parseFuncColor reads the functional forms: `rgb`, `rgba` and `hsl`, with the
// numbers separated by commas, by spaces, or by a slash before the alpha, as
// SVG 2 and CSS both write them.
func parseFuncColor(name, args string) (canvas.Color, bool) {
	// Split the argument list into numbers and percentages, dropping the
	// punctuation that separates them, keeping each one's unit.
	fields := strings.FieldsFunc(args, func(r rune) bool {
		return r == ',' || r == '/' || r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})
	if len(fields) < 3 || len(fields) > 4 {
		return 0, false
	}
	a := 255.0
	if len(fields) == 4 {
		v, ok := colorAlpha(fields[3])
		if !ok {
			return 0, false
		}
		a = float64(v)
	}
	switch strings.ToLower(name) {
	case "rgb", "rgba":
		r, ok1 := colorChannel(fields[0])
		g, ok2 := colorChannel(fields[1])
		b, ok3 := colorChannel(fields[2])
		if !ok1 || !ok2 || !ok3 {
			return 0, false
		}
		return canvas.RGBA(canvas.RGBAParams{R: toByte(r), G: toByte(g), B: toByte(b), A: toByte(a)}), true
	case "hsl", "hsla":
		h, err1 := parseHue(fields[0])
		sv, err2 := parsePercent(fields[1])
		lv, err3 := parsePercent(fields[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return 0, false
		}
		r, g, b := hslToRGB(h, sv/100, lv/100)
		return canvas.RGBA(canvas.RGBAParams{R: toByte(r * 255), G: toByte(g * 255), B: toByte(b * 255), A: toByte(a)}), true
	}
	return 0, false
}
