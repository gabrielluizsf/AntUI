package css

import (
	"strings"
)

// sourceKind reads one @font-face src source: which function it is and what
// it names, or "" when the source is neither url() nor local().
func sourceKind(part string) (kind, arg string) {
	s := strings.TrimSpace(part)
	for _, k := range []string{"url", "local"} {
		if len(s) > len(k)+1 && strings.EqualFold(s[:len(k)], k) && s[len(k)] == '(' {
			end := parenEnd(s, len(k))
			if end < 0 {
				return "", ""
			}
			return strings.ToLower(k), unquote(strings.TrimSpace(s[len(k)+1 : end]))
		}
	}
	return "", ""
}
