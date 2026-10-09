package css

import (
	"fmt"
	"strings"
)

func parseFour(raw string) ([4]Length, error) {
	var out [4]Length
	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return out, fmt.Errorf("css: expected a length, got %q", raw)
	}
	vals := make([]Length, len(parts))
	for i, p := range parts {
		v, err := parseLength(p)
		if err != nil {
			return out, err
		}
		vals[i] = v
	}
	return expandInto(out, vals), nil
}

// parseFourAt is parseFour under a measurement context, so % and the
// viewport/font units resolve and box shorthands may hold a calc(). A box
// shorthand that is a single math formula resolves to one equal length.
