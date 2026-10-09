package css

import (
	"strconv"
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

// readEscape resolves the escape sequence starting at the '\' in s. A
// backslash followed by up to six hex digits reads a code point (with one
// optional trailing whitespace as its terminator); any other backslash reads
// the literal next character.
func readEscape(s string, i int) (string, int) {
	// s[i] is '\\'.
	if i+1 >= len(s) {
		return "", 0
	}
	n := i + 1
	if !isHex(s[n]) {
		return string(s[n]), 2
	}
	start := n
	for n < len(s) && n < start+6 && isHex(s[n]) {
		n++
	}
	v, err := strconv.ParseUint(s[start:n], 16, 32)
	if err != nil {
		return "", 0
	}
	if v > 0x10FFFF {
		v = 0xFFFD
	}
	if n < len(s) && isSpace(s[n]) {
		n++ // the single optional whitespace terminator of a hex escape
	}
	return string(rune(v)), n - i
}
