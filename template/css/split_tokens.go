package css

// splitTokens splits a value on top-level whitespace, keeping a parenthesised
// group such as rgb(0 0 0) or drop-shadow(0 1px 2px black) in one piece.
func splitTokens(s string) []string {
	var out []string
	depth, start := 0, -1
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ' ', '\t', '\n', '\r':
			if depth == 0 {
				if start >= 0 {
					out = append(out, s[start:i])
					start = -1
				}
				continue
			}
		}
		if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		out = append(out, s[start:])
	}
	return out
}
