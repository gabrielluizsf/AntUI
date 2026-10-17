package svg

import "strings"

// splitDecl splits one declaration at the colon that separates its name from
// its value, and not at a colon inside brackets or a quoted string, as the
// `http:` in `url(http://…)` is.
func splitDecl(decl string) (name, value string, ok bool) {
	depth := 0
	var quote byte
	for i := 0; i < len(decl); i++ {
		c := decl[i]
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '"' || c == '\'':
			quote = c
		case c == '(' || c == '[':
			depth++
		case c == ')' || c == ']':
			depth--
		case c == ':' && depth == 0:
			name = declName(decl[:i])
			if name == "" {
				return "", "", false
			}
			return name, strings.TrimSpace(decl[i+1:]), true
		}
	}
	return "", "", false
}
