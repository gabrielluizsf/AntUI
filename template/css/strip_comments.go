package css

import (
	"strings"
)

// stripComments removes every /* … */ comment from a string. A comment in a
// file any old position — selector, value, media prelude — carries nothing, so
// it is just deleted.
func stripComments(s string) string {
	for {
		a := strings.Index(s, "/*")
		if a < 0 {
			return s
		}
		b := strings.Index(s[a+2:], "*/")
		if b < 0 {
			return s[:a]
		}
		end := a + 2 + b + 2
		s = s[:a] + s[end:]
	}
}
