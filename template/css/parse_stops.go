package css

import (
	"strings"
)

// parseStops reads the colour stops of a gradient. A bare percentage or
// length between stops is a colour hint and is skipped, the way a browser
// reads the hint while keeping the interpolation the two neighbours ask for.
func parseStops(parts []string, ctx Units) ([]Stop, bool) {
	var out []Stop
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if isStopHint(p, ctx) {
			continue
		}
		stop, ok := parseStop(p, ctx)
		if !ok {
			return nil, false
		}
		out = append(out, stop)
	}
	return out, len(out) > 0
}
