package svg

import (
	"strconv"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// parseFuncColor reads the functional forms: `rgb`, `rgba` and `hsl`, with the
// numbers separated by commas, by spaces, or by a slash before the alpha, as
// SVG 2 and CSS both write them.
func parseFuncColor(name, args string) (canvas.Color, bool) {
	// split the argument list into numbers and percentages, dropping the
	// punctuation that separates them, keeping each one's unit.
	fields := strings.FieldsFunc(args, func(r rune) bool {
		return r == ',' || r == '/' || r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})
	if len(fields) < 3 || len(fields) > 4 {
		return 0, false
	}
	// channel reads one number that is either 0-255 or a percentage of it.
	channel := func(s string) (float64, bool) {
		if p, ok := strings.CutSuffix(s, "%"); ok {
			v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
			if err != nil {
				return 0, false
			}
			return v * 255 / 100, true
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
		if err != nil {
			return 0, false
		}
		return v, true
	}
	alpha := func(s string) (uint8, bool) {
		if p, ok := strings.CutSuffix(s, "%"); ok {
			v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
			if err != nil {
				return 0, false
			}
			return byte(clampFloat(v*255/100, 0, 255) + 0.5), true
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
		if err != nil {
			return 0, false
		}
		return byte(clampFloat(v*255, 0, 255) + 0.5), true
	}
	a := 255.0
	if len(fields) == 4 {
		v, ok := alpha(fields[3])
		if !ok {
			return 0, false
		}
		a = float64(v)
	}

	switch strings.ToLower(name) {
	case "rgb", "rgba":
		r, ok1 := channel(fields[0])
		g, ok2 := channel(fields[1])
		b, ok3 := channel(fields[2])
		if !ok1 || !ok2 || !ok3 {
			return 0, false
		}
		return canvas.RGBA(toByte(r), toByte(g), toByte(b), toByte(a)), true
	case "hsl", "hsla":
		h, err1 := parseHue(fields[0])
		sv, err2 := parsePercent(fields[1])
		lv, err3 := parsePercent(fields[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return 0, false
		}
		r, g, b := hslToRGB(h, sv/100, lv/100)
		return canvas.RGBA(toByte(r*255), toByte(g*255), toByte(b*255), toByte(a)), true
	}
	return 0, false
}
