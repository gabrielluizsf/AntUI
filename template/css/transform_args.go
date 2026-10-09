package css

import (
	"strings"
)

// transformArgs splits a function's arguments on commas and whitespace and
// rejects leftovers that are not real arguments.
func transformArgs(inner string) ([]string, bool) {
	if strings.HasPrefix(inner, ",") || strings.HasSuffix(inner, ",") || strings.Contains(inner, ",,") {
		return nil, false
	}
	var out []string
	for _, f := range strings.FieldsFunc(inner, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' || r == '\n' || r == '\r'
	}) {
		if strings.ContainsAny(f, "()") {
			return nil, false
		}
		out = append(out, f)
	}
	return out, true
}
