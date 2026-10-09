package css

import "strings"

// gridFunction splits "name(arg…)" into its name and parenthesised argument list.
func gridFunction(raw string) (string, string, bool) {
	raw = strings.TrimSpace(raw)
	open := strings.IndexByte(raw, '(')
	if open <= 0 || !strings.HasSuffix(raw, ")") || parenEnd(raw, open) != len(raw)-1 {
		return "", "", false
	}
	name := strings.ToLower(strings.TrimSpace(raw[:open]))
	return name, strings.TrimSpace(raw[open+1 : len(raw)-1]), true
}
