package svg

// skipSpace passes over whitespace and nothing else.
func (p *xmlParser) skipSpace() {
	for !p.eof() && isSpace(p.src[p.pos]) {
		p.pos++
	}
}
