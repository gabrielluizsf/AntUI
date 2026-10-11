package svg

// skipDecl passes over a `<!…>` declaration, counting brackets so that the
// internal subset of a DTD, which can hold a `>` of its own, does not end it
// early.
func (p *xmlParser) skipDecl() {
	depth := 0
	for !p.eof() {
		switch p.src[p.pos] {
		case '[':
			depth++
		case ']':
			depth--
		case '>':
			if depth <= 0 {
				p.pos++
				return
			}
		}
		p.pos++
	}
}
