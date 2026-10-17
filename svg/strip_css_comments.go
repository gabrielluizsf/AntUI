package svg

import "strings"

// stripCSSComments takes the `/* … */` comments out of a sheet: they are not
// rules and hold nothing a drawing reads.
func stripCSSComments(s string) string {
	for {
		i := strings.Index(s, "/*")
		if i < 0 {
			return s
		}
		j := strings.Index(s[i:], "*/")
		if j < 0 {
			return s[:i]
		}
		s = s[:i] + " " + s[i+j+2:]
	}
}

// blockEnd answers where the block opened at open is closed, following the
// braces nested inside it so that an at-rule holding rules of its own is one
// block and not the first of them, or -1 where it is never closed.
func blockEnd(s string, open int) int {
	depth := 0
	for i := open; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}
