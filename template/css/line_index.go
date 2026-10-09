package css

// LineIndex finds the line carrying a name: the first one counting from the
// start, the last one counting from the end — the way a placement resolves its
// start line and its end line. Implicit lines hold no names, so a name no
// explicit line carries is simply not found.
func (a GridAxis) LineIndex(name string, fromEnd bool) (int, bool) {
	if fromEnd {
		for i := len(a.Names) - 1; i >= 0; i-- {
			if gridHasName(a.Names[i], name) {
				return i, true
			}
		}
		return 0, false
	}
	for i, names := range a.Names {
		if gridHasName(names, name) {
			return i, true
		}
	}
	return 0, false
}
