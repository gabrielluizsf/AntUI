package svg

import "fmt"

// inQuotes reads past a quoted part of a declaration, where a `>` or a bracket
// inside it is a character of the string and not the end of the declaration.
func (p *xmlParser) inQuotes() error {
	quote := p.src[p.pos]
	p.pos++
	for !p.eof() {
		if p.src[p.pos] == quote {
			p.pos++
			return nil
		}
		p.pos++
	}
	return fmt.Errorf("antui/svg: a quoted part of a declaration is never closed")
}
