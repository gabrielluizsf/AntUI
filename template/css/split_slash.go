package css

// splitSlash splits on the "/" that sits at the top level of a colour
// function, the divider between channels and alpha. A slash inside nested
// parentheses belongs to something else and is ignored.
func splitSlash(s string) (head, tail string, has bool) {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
		case '/':
			if depth == 0 {
				return s[:i], s[i+1:], true
			}
		}
	}
	return s, "", false
}
