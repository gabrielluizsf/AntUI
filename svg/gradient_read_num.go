package svg

import (
	"strconv"
	"strings"
)

// readNum reads one named numeric attribute of a gradient. A percentage is
// kept as it was written, because which of the two things it is a percentage
// of is decided by the units, and the units may still be coming from an href
// that has not been followed yet — only that it is a number is worth saying
// now. A length is turned into a number; anything else is warned about and
// left at what it was.
func (g *gradient) readNum(attrs map[string]string, name string, dst *float64, axis int, warn func(string, ...any)) {
	raw, ok := attrs[name]
	if !ok {
		return
	}
	g.set[name] = true
	trimmed := strings.TrimSpace(raw)
	if body, isPct := strings.CutSuffix(trimmed, "%"); isPct {
		if _, err := strconv.ParseFloat(strings.TrimSpace(body), 64); err != nil {
			warn("%s %q is not a number, so it stays at %g", name, raw, *dst)
			return
		}
		g.pct[name] = pctRead{raw: trimmed, axis: axis}
		return
	}
	v, ok := gradientLength(raw, 0)
	if !ok {
		warn("%s %q is not a length, so it stays at %g", name, raw, *dst)
		return
	}
	*dst = v
}
