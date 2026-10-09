package css

// closedGridString reports whether every quote in a row opens and closes before the row ends.
func closedGridString(raw string) bool {
	if len(raw) < 2 || (raw[0] != '"' && raw[0] != '\'') {
		return false
	}
	quote := raw[0]
	for i := 1; i < len(raw)-1; i++ {
		if raw[i] == '\\' {
			i++
			continue
		}
		if raw[i] == quote {
			return false
		}
	}
	return true
}
