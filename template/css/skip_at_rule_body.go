package css

import (
	"strings"
)

// skipAtRuleBody swallows the rest of an unknown at-rule: a ';' statement or
// a balanced { … } block, honouring strings and comments.
func (p *parser) skipAtRuleBody() error {
	for !p.eof() {
		c := p.peek()
		switch {
		case c == '"' || c == '\'':
			if end := quoteEnd(p.src, p.i); end >= 0 {
				p.i = end + 1
			} else {
				p.i = len(p.src)
			}
		case c == '/':
			if p.i+1 < len(p.src) && p.src[p.i+1] == '*' {
				if end := strings.Index(p.src[p.i+2:], "*/"); end >= 0 {
					p.i += end + 4
					continue
				}
				p.i = len(p.src)
				return nil
			}
			p.next()
		case c == ';':
			p.next()
			return nil
		case c == '{':
			return p.skipBlock()
		case c == '}':
			return nil
		default:
			p.next()
		}
	}
	return nil
}
