package css

import (
	"strings"
)

// parseDirectionList reads a comma-separated list of direction keywords.
func parseDirectionList(raw string) ([]uint8, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	var out []uint8
	for _, p := range splitFields(raw, ',') {
		v, ok := parseDirection(strings.TrimSpace(p))
		if !ok {
			return nil, false
		}
		out = append(out, v)
	}
	return out, true
}
