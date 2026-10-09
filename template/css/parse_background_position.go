package css

import (
	"strings"
)

// parseBackgroundPosition reads a comma-separated background-position list.
// Each value holds one position per axis; a single value is expanded across
// the other axis the way CSS does, centring it. A malformed value drops the
// whole property.
func parseBackgroundPosition(raw string, ctx Units) ([]BackPos, bool) {
	var out []BackPos
	for _, part := range splitFields(raw, ',') {
		p, ok := parseBackPosTokens(splitTokens(strings.TrimSpace(part)), ctx)
		if !ok {
			return nil, false
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}
