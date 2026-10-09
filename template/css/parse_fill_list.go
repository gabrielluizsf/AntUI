package css

import (
	"strings"
)

// parseFillList reads a comma-separated list of fill-mode keywords.
func parseFillList(raw string) ([]uint8, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	var out []uint8
	for _, p := range splitFields(raw, ',') {
		v, ok := parseFill(strings.TrimSpace(p))
		if !ok {
			return nil, false
		}
		out = append(out, v)
	}
	return out, true
}
