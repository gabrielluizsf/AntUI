package svg

// validName reports whether s is a name XML allows: one that opens with a
// letter, an underscore or a colon — never a digit, a quote or a bracket — and
// holds only letters, digits, dashes, dots, underscores, colons and the bytes
// of a written-out character. A tag or an attribute with any other spelling is
// not a name at all, and a drawing that writes one is not the drawing it
// declares. Bytes at or above 0x80 are let through untouched, because a name
// may be written in any script and this is not the place to argue with one.
func validName(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 0x80:
			continue
		case i == 0:
			if !isNameStart(c) {
				return false
			}
		case isNameChar(c):
			continue
		default:
			return false
		}
	}
	return true
}

// isNameStart is the set a name may open with: a letter, an underscore or a
// colon.
func isNameStart(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_' || c == ':'
}

// isNameChar is the set a name may hold after its first byte: everything it
// may open with, plus the digits, dashes and dots that only appear within one.
func isNameChar(c byte) bool {
	return isNameStart(c) || '0' <= c && c <= '9' || c == '-' || c == '.'
}
