package css

import (
	"strings"
)

// readDecl reads one declaration inside a block, stopping at the ';' or the
// '}' that ends it. Strings, comments and parentheses are kept intact, so a
// value like url("a;b.png") reads as one unit.
func (p *parser) readDecl() (string, error) {
	start := p.i
	depth := 0
	for !p.eof() {
		c := p.src[p.i]
		switch {
		case c == '/' && p.i+1 < len(p.src) && p.src[p.i+1] == '*':
			if end := strings.Index(p.src[p.i+2:], "*/"); end >= 0 {
				p.i += end + 4
			} else {
				p.i = len(p.src)
			}
		case c == '"' || c == '\'':
			if end := quoteEnd(p.src, p.i); end >= 0 {
				p.i = end + 1
			} else {
				return "", fmtErrf("unterminated string")
			}
		case c == '(':
			depth++
			p.i++
		case c == ')':
			if depth > 0 {
				depth--
			}
			p.i++
		case depth == 0 && (c == ';' || c == '}'):
			return p.src[start:p.i], nil
		default:
			p.i++
		}
	}
	return "", fmtErrf("unterminated declaration")
}
