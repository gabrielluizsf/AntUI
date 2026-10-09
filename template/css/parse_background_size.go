package css

import (
	"strings"
)

// parseBackgroundSize reads a comma-separated background-size list: each value
// is "cover", "contain" or one or two lengths, the missing axis becoming auto.
func parseBackgroundSize(raw string, ctx Units) ([]BackSize, bool) {
	var out []BackSize
	for _, part := range splitFields(raw, ',') {
		sz, ok := parseBackSizeOne(strings.TrimSpace(part), ctx)
		if !ok {
			return nil, false
		}
		out = append(out, sz)
	}
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}
