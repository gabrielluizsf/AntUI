package css

import "github.com/gabrielluizsf/antui/canvas"

func parseBorderColors(raw string) ([4]canvas.Color, bool) {
	var out [4]canvas.Color
	parts := splitWords(raw)
	if len(parts) == 0 || len(parts) > 4 {
		return out, false
	}
	vals := make([]canvas.Color, len(parts))
	for i, p := range parts {
		c, err := ParseColor(p)
		if err != nil {
			return out, false
		}
		vals[i] = c
	}
	fill := expandFourColor(vals)
	copy(out[:], fill[:])
	return out, true
}
