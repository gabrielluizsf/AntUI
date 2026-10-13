package css

import (
	"fmt"
	"strings"
)

func parseFourAt(raw string, ctx Units) ([sideCount]Length, error) {
	if mathLook(raw) {
		px, ok := EvalMath(raw, ctx)
		if ok {
			return [sideCount]Length{Fixed(px), Fixed(px), Fixed(px), Fixed(px)}, nil
		}
		return [sideCount]Length{}, fmt.Errorf("css: %q does not evaluate", raw)
	}
	var out [sideCount]Length
	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return out, fmt.Errorf("css: expected a length, got %q", raw)
	}
	vals := make([]Length, len(parts))
	for i, p := range parts {
		v, err := parseLengthAt(p, ctx)
		if err != nil {
			return out, err
		}
		vals[i] = v
	}
	return expandFour(vals), nil
}
