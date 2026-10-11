package svg

// attrName reads the name of one attribute, stopping at the equals sign or the
// space after it.
func (p *xmlParser) attrName() (string, bool) {
	start := p.pos
	for !p.eof() {
		c := p.src[p.pos]
		if c == '=' || isSpace(c) || c == '>' || c == '/' {
			break
		}
		p.pos++
	}
	if p.pos == start {
		return "", false
	}
	return p.src[start:p.pos], true
}
