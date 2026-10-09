package css

import (
	"strings"
)

// readHeader reads a prelude (a selector list or an at-rule condition) up to
// the '{' that opens the block, returning the text it passed over. When semi
// is true, a top-level ';' also ends the header, which a stray statement uses
// to drop itself. Strings, /* comments */ and parentheses are kept intact.
func (p *parser) readHeader(semi bool) (string, bool, error) {
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
				return "", false, fmtErrf("unterminated string")
			}
		case c == '(':
			depth++
			p.i++
		case c == ')':
			if depth > 0 {
				depth--
			}
			p.i++
		case depth == 0 && c == '{':
			return p.src[start:p.i], false, nil
		case depth == 0 && semi && c == ';':
			return p.src[start:p.i], true, nil
		default:
			p.i++
		}
	}
	return "", false, fmtErrf("unterminated rule, missing '{'")
}
