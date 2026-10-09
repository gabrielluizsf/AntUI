package css

import (
	"strings"
)

// skipSpace advances past whitespace and /* comments */.
func (p *parser) skipSpace() {
	for !p.eof() {
		c := p.src[p.i]
		if c == '/' && p.i+1 < len(p.src) && p.src[p.i+1] == '*' {
			end := strings.Index(p.src[p.i+2:], "*/")
			if end < 0 {
				p.i = len(p.src)
				return
			}
			p.i += end + 4
			continue
		}
		if isSpace(c) {
			p.i++
			continue
		}
		return
	}
}
