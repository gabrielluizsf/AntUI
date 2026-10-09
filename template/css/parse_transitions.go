package css

import (
	"strings"
)

// parseTransitions reads the transition shorthand: a comma-separated list of
// transitions. "none" clears every entry, so a style returns no transitions
// at all. An empty or unrecognised entry invalidates the whole property, the
// way a browser drops a declaration it cannot parse.
func parseTransitions(raw string) ([]Transition, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	if raw == "none" {
		return nil, true
	}
	var out []Transition
	for _, part := range splitFields(raw, ',') {
		tr, ok := parseTransition(strings.TrimSpace(part))
		if !ok {
			return nil, false
		}
		out = append(out, tr)
	}
	return out, true
}
