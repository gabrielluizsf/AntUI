package css

import "strings"

// splitGridSlashes splits a grid value on the "/" that separates rows from columns, ignoring slashes inside a group.
func splitGridSlashes(raw string) []string {
	if raw == "" {
		return nil
	}
	var out []string
	start, depth := -1, 0
	hadSlash := false
	for i := 0; i < len(raw); i++ {
		switch raw[i] {
		case '(':
			if start < 0 {
				start = i
			}
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case '/':
			if depth == 0 {
				if start < 0 {
					start = i
				}
				out = append(out, strings.TrimSpace(raw[start:i]))
				start = -1
				hadSlash = true
			}
		default:
			if start < 0 {
				start = i
			}
		}
	}
	if start >= 0 {
		out = append(out, strings.TrimSpace(raw[start:]))
	} else if hadSlash {
		out = append(out, "")
	}
	return out
}
