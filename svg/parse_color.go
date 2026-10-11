package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// parseColor reads a colour in any of the ways a drawing may write one: the hex
// forms from three to eight digits, the functional forms with numbers or
// percentages, and the name of one of the colours CSS knows. The bool says
// whether it was one this could read, so a caller can carry on with the rest of
// the drawing rather than stopping at a colour it has never heard of.
func parseColor(s string) (canvas.Color, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	if s[0] == '#' {
		return parseHexColor(s[1:])
	}
	if name, args, ok := strings.Cut(s, "("); ok && strings.HasSuffix(args, ")") {
		return parseFuncColor(name, args[:len(args)-1])
	}
	if c, ok := canvas.ParseColor(s); ok == nil {
		return c, true
	}
	if c, ok := namedColors[strings.ToLower(s)]; ok {
		return c, true
	}
	return 0, false
}
