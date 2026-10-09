package css

// splitQueryTokens breaks a media query into its words and (…)-groups,
// keeping parenthesised groups — including nested parentheses — whole.
func splitQueryTokens(s string) []string {
	var out []string
	i := 0
	for i < len(s) {
		for i < len(s) && isSpace(s[i]) {
			i++
		}
		if i >= len(s) {
			break
		}
		start := i
		if s[i] == '(' {
			i = max(scanParen(s, i), i+1)
		} else {
			for i < len(s) && !isSpace(s[i]) && s[i] != '(' {
				i++
			}
		}
		if i == start {
			i++
			continue
		}
		out = append(out, s[start:i])
	}
	return out
}
