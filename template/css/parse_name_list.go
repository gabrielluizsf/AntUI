package css

import (
	"strings"
)

// parseNameList reads an animation-name value: a comma-separated series of
// keyframe names.
func parseNameList(raw string) ([]string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	var out []string
	for _, p := range splitFields(raw, ',') {
		p = trimQuotes(strings.TrimSpace(p))
		if p == "" {
			return nil, false
		}
		out = append(out, p)
	}
	return out, true
}
