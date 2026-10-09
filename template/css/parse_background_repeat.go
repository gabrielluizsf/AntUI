package css

import (
	"strings"
)

// parseBackgroundRepeat reads a comma-separated background-repeat list: one
// keyword for both axes, or an X and a Y keyword, with the shorthand
// repeat-x and repeat-y folded into their pair.
func parseBackgroundRepeat(raw string) ([]BackRepeat, bool) {
	var out []BackRepeat
	for _, part := range splitFields(raw, ',') {
		rep, ok := backRepeatPair(splitTokens(strings.TrimSpace(part)))
		if !ok {
			return nil, false
		}
		out = append(out, rep)
	}
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}
