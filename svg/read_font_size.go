package svg

import (
	"strconv"
	"strings"
)

// readFontSize is how tall the writing is: a length in the drawing's own units,
// a percentage of the size it already has, one of the sizes CSS gives a name
// to, or a step up or down from the size it already has. Anything else leaves
// that size alone, since a size this cannot read is not a size to write at.
func readFontSize(was float64, raw string, warn func(string, ...any)) float64 {
	if was <= 0 {
		was = defaultFontSize
	}
	s := strings.ToLower(strings.TrimSpace(raw))
	switch s {
	case "larger":
		return was * 1.2
	case "smaller":
		return was / 1.2
	}
	if named, ok := namedFontSizes[s]; ok {
		return named
	}
	if p, ok := strings.CutSuffix(s, "%"); ok {
		if v, err := strconv.ParseFloat(strings.TrimSpace(p), 64); err == nil && v > 0 {
			return was * v / 100
		}
		warn("the font size %q is not a percentage, so the writing keeps the size it had", raw)
		return was
	}
	if v, ok := parseLength(s); ok && v > 0 {
		return v
	}
	warn("the font size %q is not a length, so the writing keeps the size it had", raw)
	return was
}
