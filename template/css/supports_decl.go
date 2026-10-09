package css

import (
	"strings"
)

// splitSupportDecl reads "property: value" when a colon at the top level of
// the group splits it there. Inside parens or quotes a colon belongs to what
// it is written in — url(http://x), "a:b", selector(a:hover) — and a group
// whose first colon is one of those is not a declaration at all.
func splitSupportDecl(text string) (Declaration, bool) {
	top := -1
	depth := 0
	for i := 0; i < len(text); i++ {
		c := text[i]
		switch {
		case c == '"' || c == '\'':
			end := quoteEnd(text, i)
			if end < 0 {
				return Declaration{}, false
			}
			i = end
		case c == '(':
			depth++
		case c == ')':
			if depth > 0 {
				depth--
			}
		case c == ':' && depth == 0:
			top = i
		}
		if top >= 0 {
			break
		}
	}
	if top < 0 || strings.IndexByte(text, ':') != top {
		return Declaration{}, false
	}
	return parseDeclaration(text)
}
