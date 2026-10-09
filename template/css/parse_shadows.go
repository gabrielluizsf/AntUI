package css

import (
	"strings"
)

// parseShadows reads a comma-separated shadow list. allowInset admits the
// inset keyword (box-shadow) and max is the number of lengths the property
// takes: four for box-shadow, three for text-shadow. "none" clears the list.
func parseShadows(raw string, ctx Units, allowInset bool, max int) ([]Shadow, bool) {
	if raw == "none" {
		return nil, true
	}
	var out []Shadow
	for _, part := range splitFields(raw, ',') {
		if sh, ok := parseShadow(strings.TrimSpace(part), ctx, allowInset, max); ok {
			out = append(out, sh)
		}
	}
	return out, len(out) > 0
}
