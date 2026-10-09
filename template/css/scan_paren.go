package css

// scanParen returns the index just past the ')' matching the '(' at i, or the
// end of the string when the group never closes.
func scanParen(s string, i int) int {
	end := parenEnd(s, i)
	if end < 0 {
		return len(s)
	}
	return end + 1
}
