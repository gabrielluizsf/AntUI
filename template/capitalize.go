package template

import (
	"strings"
	"unicode"
)

// capitalize uppercases the first letter of every word.
func capitalize(s string) string {
	var b strings.Builder
	start := true
	for _, r := range s {
		if unicode.IsLetter(r) && start {
			b.WriteRune(unicode.ToUpper(r))
			start = false
			continue
		}
		start = !unicode.IsLetter(r) && !unicode.IsDigit(r)
		b.WriteRune(r)
	}
	return b.String()
}
