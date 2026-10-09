package css

import (
	"strconv"
	"strings"
)

// parseLineHeight reads a unitless multiplier, a percentage, or a length,
// always computing a multiplier of the font size. Zero means normal.
func parseLineHeight(raw string, ctx Units) (float64, bool) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return 0, false
	}
	if s == "normal" {
		return 0, true
	}
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil || v < 0 {
			return 0, false
		}
		return v / 100, true
	}
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		if v < 0 {
			return 0, false
		}
		return v, true
	}
	l, err := parseLengthAt(s, ctx)
	if err != nil {
		return 0, false
	}
	return float64(l.Resolve(ctx)) / float64(ctx.font()), true
}
