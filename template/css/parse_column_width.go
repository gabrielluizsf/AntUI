package css

import (
	"strings"
)

// parseColumnWidth reads the width a column is given when the container is
// wide enough for it. A length is a width; auto leaves the count to decide.
func parseColumnWidth(raw string, ctx Units) (Length, bool) {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "auto" {
		return Auto(), true
	}
	l, err := parseLengthAt(s, ctx)
	if err != nil || l.Auto() || l.None() || l.value < 0 {
		return Auto(), false
	}
	return l, true
}
