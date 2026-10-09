package css

import (
	"strings"
)

// parseAnimations reads the animation shorthand: a comma-separated list of
// animations, each made of a keyframe name, two times, an easing, an
// iteration count, and the direction and fill keywords in any order. An entry
// with no name — an empty list, a stray keyword — invalidates the whole
// property, exactly as a browser drops an unparsable declaration.
func parseAnimations(raw string) ([]Animation, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	var out []Animation
	for _, part := range splitFields(raw, ',') {
		a, ok := parseAnimation(strings.TrimSpace(part))
		if !ok {
			return nil, false
		}
		out = append(out, a)
	}
	return out, len(out) > 0
}
