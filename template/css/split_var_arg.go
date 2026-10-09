package css

import "strings"

// splitVarArg reads " --name , fallback " into the custom-property name, the
// fallback (may be empty) and ok. The fallback is everything after the first
// top-level comma.
func splitVarArg(inner string) (name, fallback string, ok bool) {
	inside := false
	for i := 0; i < len(inner); i++ {
		c := inner[i]
		switch c {
		case '"', '\'':
			if j := quoteEnd(inner, i); j >= 0 {
				i = j
			}
		case '(', '[':
			inside = true
		case ')', ']':
			inside = false
		case ',':
			if !inside {
				name = strings.TrimSpace(inner[:i])
				fallback = strings.TrimSpace(inner[i+1:])
				return name, fallback, name != ""
			}
		}
	}
	return strings.TrimSpace(inner), "", strings.TrimSpace(inner) != ""
}
