package css

import (
	"strings"
)

// mathLook reports whether a raw value opens a math formula the engine can
// evaluate (calc/min/max/clamp). Anything else with a "(" in the way is a
// function this engine does not understand, not something to attempt.
func mathLook(raw string) bool {
	t := strings.ToLower(strings.TrimSpace(stripComments(raw)))
	for _, fn := range []string{"calc(", "min(", "max(", "clamp("} {
		if strings.HasPrefix(t, fn) {
			return true
		}
	}
	return false
}
