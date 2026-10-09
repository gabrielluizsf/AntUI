package css

import (
	"strings"
)

// parseGridLineNames reads a "[…]" line-name list. Names are case-sensitive,
// the way CSS spells them, and the keywords the grid grammar already spent are
// refused: a line named "auto" or "span" would read as a placement, not a name.
func parseGridLineNames(raw string) ([]string, bool) {
	raw = strings.TrimSpace(raw)
	if len(raw) < 3 || raw[0] != '[' || raw[len(raw)-1] != ']' {
		return nil, false
	}
	fields := strings.Fields(raw[1 : len(raw)-1])
	if len(fields) == 0 {
		return nil, false
	}
	for _, field := range fields {
		if !validGridName(field) {
			return nil, false
		}
	}
	return fields, true
}
