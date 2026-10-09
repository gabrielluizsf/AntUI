package css

// isUnitByte reports whether a byte can continue a numeric literal's unit.
func isUnitByte(c byte) bool {
	switch {
	case c >= 'a' && c <= 'z':
		return true
	case c >= 'A' && c <= 'Z':
		return true
	case c == '%':
		return true
	}
	return false
}
