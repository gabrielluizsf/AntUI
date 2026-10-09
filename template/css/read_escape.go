package css

import "strconv"

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
