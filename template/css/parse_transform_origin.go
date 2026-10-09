package css

import (
	"strings"
)

// parseTransformOrigin reads 1–3 values: a vertical keyword alone sets x to
// the centre; two values order left/right/centre for x and top/centre/bottom
// for y in either spelling; a third, the z, is validated and dropped.
func parseTransformOrigin(raw string, ctx Units) ([2]Length, bool) {
	parts := strings.Fields(strings.ToLower(raw))
	if len(parts) == 0 || len(parts) > 3 {
		return [2]Length{}, false
	}
	if len(parts) == 1 {
		if y, ok := originKeywordY(parts[0]); ok {
			return [2]Length{Pct(50), y}, true
		}
		x, ok := originKeywordX(parts[0], ctx)
		if !ok {
			return [2]Length{}, false
		}
		return [2]Length{x, Pct(50)}, true
	}
	xTok, yTok := parts[0], parts[1]
	if parts[0] == "top" || parts[0] == "bottom" {
		xTok, yTok = parts[1], parts[0]
	}
	x, ok := originKeywordX(xTok, ctx)
	if !ok {
		return [2]Length{}, false
	}
	y, ok := originKeywordYWithX(yTok, ctx)
	if !ok {
		return [2]Length{}, false
	}
	if len(parts) == 3 {
		if _, err := parseLengthAt(parts[2], ctx); err != nil {
			if _, ok := transformNumber(parts[2]); !ok {
				return [2]Length{}, false
			}
		}
	}
	return [2]Length{x, y}, true
}
