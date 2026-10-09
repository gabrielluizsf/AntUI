package css

// gridHasName reports whether one line carries a name.
func gridHasName(names []string, name string) bool {
	for _, candidate := range names {
		if candidate == name {
			return true
		}
	}
	return false
}
