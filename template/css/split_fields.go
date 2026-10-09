package css

// splitFields splits a comma-joined list, ignoring commas nested in
// parentheses.
func splitFields(s string, sep byte) []string {
	var out []string
	start, depth := 0, 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
		case sep:
			if depth == 0 {
				out = append(out, s[start:i])
				start = i + 1
			}
		}
	}
	out = append(out, s[start:])
	return out
}
