package css

import (
	"strings"
)

// skipBlock swallows a balanced { … } block.
func (p *parser) skipBlock() error {
	depth := 0
	for !p.eof() {
		c := p.src[p.i]
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
		case c == '{':
			depth++
			p.next()
		case c == '}':
			depth--
			p.next()
			if depth == 0 {
				return nil
			}
		default:
			p.next()
		}
	}
	return fmtErrf("unterminated block")
}

// parseSelectors splits a comma-separated selector list into one Selector per
