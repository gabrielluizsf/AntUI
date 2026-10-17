package svg

import "fmt"

// entityValue reads the quoted text a `<!ENTITY …>` stands for. A name the
// doctype leans on itself — %name; — is read into it here, since the doctype
// is read once where the drawing's text is read every time.
func (p *xmlParser) entityValue(to int) (string, error) {
	if p.pos >= to || !p.at("\"") && !p.at("'") {
		return "", fmt.Errorf("antui/svg: its text was expected in quotes")
	}
	quote := p.src[p.pos]
	p.pos++
	start := p.pos
	for p.pos < to {
		if p.src[p.pos] == quote {
			text := p.src[start:p.pos]
			p.pos++
			return p.expandPE(text)
		}
		p.pos++
	}
	return "", fmt.Errorf("antui/svg: the text of an entity is never closed")
}

// readName reads the name a declaration gives, and answers it — or nothing,
// when no name is written where one was expected. Bytes at or above 0x80 are
// let through as the rest of the name they are part of, since a name may be
// written in any script.
func (p *xmlParser) readName() string {
	start := p.pos
	for !p.eof() {
		c := p.src[p.pos]
		if p.pos == start {
			if !isNameStart(c) && c < 0x80 {
				break
			}
		} else if !isNameChar(c) && c < 0x80 {
			break
		}
		p.pos++
	}
	return p.src[start:p.pos]
}
