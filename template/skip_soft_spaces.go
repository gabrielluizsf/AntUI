package template

// skipSoftSpaces advances past the soft spaces at the start of a line, which a
// collapsing white-space value drops.
func skipSoftSpaces(text string, i int) int {
	for i < len(text) && isSoftSpace(text[i]) {
		i++
	}
	return i
}
