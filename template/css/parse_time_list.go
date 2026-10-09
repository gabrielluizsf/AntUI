package css

import (
	"strings"
)

// parseTimeList reads a comma-separated list of CSS times.
func parseTimeList(raw string) ([]Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	var out []Time
	for _, p := range splitFields(raw, ',') {
		t, ok := ParseTime(strings.TrimSpace(p))
		if !ok {
			return nil, false
		}
		out = append(out, t)
	}
	return out, true
}
