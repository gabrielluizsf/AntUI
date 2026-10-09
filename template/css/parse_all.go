package css

// parseAll reads the whole of one stylesheet's text: rules and at-rules, one
// statement at a time, until the text runs out. m is the media the text sits
// inside, which for a file read on its own is none at all.
func (p *parser) parseAll(sh *Sheet, m Media) error {
	for {
		p.skipSpace()
		if p.eof() {
			return nil
		}
		switch {
		case p.peek() == '@':
			if e := p.parseAtRule(sh, m); e != nil {
				return e
			}
		case p.peek() == ';':
			p.next()
		default:
			if e := p.parseRule(sh, m); e != nil {
				return e
			}
		}
	}
}
