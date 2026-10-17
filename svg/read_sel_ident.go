package svg

// readSelIdent reads a selector name — a type, a class or an id — and answers
// it and where it ends.
func readSelIdent(s string, i int) (string, int) {
	start := i
	for i < len(s) && isNameByte(s[i], i == start) {
		i++
	}
	return s[start:i], i
}

// isNameByte says whether b may be part of a selector name. A first byte has
// to be a letter, `_` or a byte above ASCII, as CSS names have it; the rest may
// also be a digit or a dash.
func isNameByte(b byte, first bool) bool {
	if 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '_' || b >= 0x80 {
		return true
	}
	return !first && ('0' <= b && b <= '9' || b == '-')
}
