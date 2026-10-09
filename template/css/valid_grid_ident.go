package css

// validGridIdent reports whether s is a usable grid-area name.
func validGridIdent(raw string) bool {
	if raw == "" {
		return false
	}
	for i := 0; i < len(raw); i++ {
		ch := raw[i]
		if ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' || ch == '-' || ch == '_' || ch >= 0x80 {
			continue
		}
		if ch == '\\' && i+1 < len(raw) {
			i++
			continue
		}
		return false
	}
	return true
}
