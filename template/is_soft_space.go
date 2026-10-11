package template

// isSoftSpace reports whether c is a space the white-space handling manages.
func isSoftSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\f'
}
