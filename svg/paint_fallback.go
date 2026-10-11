package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// paintFallback is the colour after a `url(#g) red`, which is what the shape is
// painted with when the reference is not in the drawing. It answers whether
// there was one to give: `url(#g)` with nothing after it has none, and a
// fallback spelled `currentColor` comes back as the keyword rather than as a
// colour, because the colour of the widget is only known when it is painted.
func paintFallback(s string) (c canvas.Color, current, ok bool) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(strings.ToLower(s), "url(") {
		return 0, false, false
	}
	end := strings.Index(s, ")")
	if end < 0 {
		return 0, false, false
	}
	rest := strings.TrimSpace(s[end+1:])
	if isCurrentColor(rest) {
		return 0, true, true
	}
	c, ok = parseColor(rest)
	return c, false, ok
}
