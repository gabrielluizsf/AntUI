package css

import (
	"strings"
)

// parsePropertyList reads a transition-property value: a comma-separated
// series of canonical names, or "none" meaning nothing transitions. The
// boolean reports the none flavour; a real list reports false.
func parsePropertyList(raw string) ([]string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	if raw == "none" {
		return nil, true
	}
	var out []string
	for _, p := range splitFields(raw, ',') {
		p = strings.ToLower(strings.TrimSpace(p))
		if p == "" {
			return nil, false
		}
		out = append(out, p)
	}
	return out, false
}
