package svg

// isSpace reports whether c is one of the characters XML counts as whitespace.
func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
