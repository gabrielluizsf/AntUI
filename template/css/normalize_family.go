package css

import (
	"strings"
)

// normalizeFamily turns a font-family declaration into the name a
// @font-face answers to: comma-separated, each name unquoted, lowercased —
// a family's letters are its letters whatever case the stylesheet wrote.
func normalizeFamily(raw string) string {
	var out []string
	for _, name := range splitTopLevel(raw) {
		if n := strings.ToLower(unquote(name)); n != "" {
			out = append(out, n)
		}
	}
	return strings.Join(out, ", ")
}
