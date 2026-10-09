package css

// parenEnd returns the index of the ')' matching the '(' at open, or -1.
func parenEnd(s string, open int) int {
	depth := 0
	for i := open; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		case '"', '\'':
			if j := quoteEnd(s, i); j >= 0 {
				i = j
			}
		}
	}
	return -1
}

// quoteEnd returns the index of the quote closing the one at i, or -1.
func quoteEnd(s string, i int) int {
	q := s[i]
	for j := i + 1; j < len(s); j++ {
		if s[j] == '\\' {
			j++
			continue
		}
		if s[j] == q {
			return j
		}
	}
	return -1
}
