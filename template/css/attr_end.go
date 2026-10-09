package css

// attrEnd returns the index just past the ']' closing the attribute selector
// that starts at the '[' in s, or the end of the string when it never closes,
// in which case the selector parses to "never matches" and warns.
func attrEnd(s string, i int) int {
	for i < len(s) {
		switch s[i] {
		case '"', '\'':
			if end := quoteEnd(s, i); end >= 0 {
				i = end + 1
				continue
			}
			return len(s)
		case ']':
			return i + 1
		}
		i++
	}
	return len(s)
}
