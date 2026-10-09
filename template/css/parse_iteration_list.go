package css

import (
	"strings"
)

// parseIterationList reads a comma-separated list of iteration counts.
func parseIterationList(raw string) ([]float64, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	var out []float64
	for _, p := range splitFields(raw, ',') {
		v, ok := parseIterations(strings.TrimSpace(p))
		if !ok {
			return nil, false
		}
		out = append(out, v)
	}
	return out, true
}
