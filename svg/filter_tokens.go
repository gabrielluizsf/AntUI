package svg

// filterTokens breaks a filter list into its functions, keeping the arguments
// of each one with it: `drop-shadow(0 1px 2px black)` is one thing to try to
// read and not five, because a space inside the parentheses is part of what it
// says. Only a space at the top level, where nothing is open, comes between
// two functions.
func filterTokens(s string) []string {
	var out []string
	depth, start := 0, 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == '(':
			depth++
		case c == ')':
			if depth > 0 {
				depth--
			}
		case depth == 0 && (c == ' ' || c == '\t' || c == '\n' || c == '\r'):
			if start < i {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}
