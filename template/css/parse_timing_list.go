package css

import (
	"strings"
)

// parseTimingList reads a comma-separated list of easings.
func parseTimingList(raw string) ([]Timing, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	var out []Timing
	for _, p := range splitFields(raw, ',') {
		t, ok := parseTiming(strings.TrimSpace(p))
		if !ok {
			return nil, false
		}
		out = append(out, t)
	}
	return out, true
}
