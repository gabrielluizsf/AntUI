package css

// unhexNibble reads one hex digit into a 0-15 nibble.
func unhexNibble(c byte) (byte, bool) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', true
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, true
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}

// hexPair reads two hex digits into one byte.
func hexPair(a, b byte) (byte, bool) {
	ha, ok := unhexNibble(a)
	if !ok {
		return 0, false
	}
	hb, ok := unhexNibble(b)
	if !ok {
		return 0, false
	}
	return ha<<4 | hb, true
}
