package css

import (
	"strings"
)

// parseStop reads "red", "red 50%" or "red 10% 20%": a colour, then up to two
// offsets. Two offsets pin the colour over a span, so the painter spreads it
// from the earlier one.
func parseStop(p string, ctx Units) (Stop, bool) {
	tokens := splitTokens(p)
	if len(tokens) == 0 {
		return Stop{}, false
	}
	var first, second float64 = -1, -1
	for n := len(tokens); n > 0; n-- {
		off, err := parseStopOffset(tokens[n-1], ctx)
		if err != nil {
			break
		}
		if second < 0 {
			second = off
		} else {
			first = off
		}
		tokens = tokens[:n-1]
	}
	c, err := ParseColor(strings.Join(tokens, " "))
	if err != nil {
		return Stop{}, false
	}
	off := first
	if second >= 0 && (first < 0 || second < first) {
		off = second
	}
	return Stop{Offset: off, Color: c}, true
}
