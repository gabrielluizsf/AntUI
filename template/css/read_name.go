package css

import (
	"strings"
)

// readName reads a CSS identifier from s: ASCII letters, digits, underscore,
// hyphen, high bytes, or an escape sequence that resolves to one character.
// It returns the parsed name with escapes resolved, and the bytes consumed.
func readName(s string) (string, int) {
	var b strings.Builder
	i := 0
	for i < len(s) {
		c := s[i]
		switch {
		case c == '_' || c == '-' || isIdentByte(c) || c >= 0x80:
			b.WriteByte(c)
			i++
		case c == '\\':
			r, n := readEscape(s, i)
			if n == 0 {
				return b.String(), i
			}
			b.WriteString(r)
			i += n
		default:
			return b.String(), i
		}
	}
	return b.String(), i
}
