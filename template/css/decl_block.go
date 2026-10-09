package css

// declBlock reads the declarations inside a '{' … '}' block, one at a time,
// handing each to fn until the closing '}' — which the caller has not been
// given yet. It is the body both a rule and a @font-face share.
func (p *parser) declBlock(fn func(text string) error) error {
	for {
		p.skipSpace()
		if p.eof() {
			return fmtErrf("unterminated rule, missing '}'")
		}
		if p.peek() == '}' {
			p.next()
			return nil
		}
		text, err := p.readDecl()
		if err != nil {
			return err
		}
		if p.peek() == ';' {
			p.next()
		}
		if err := fn(text); err != nil {
			return err
		}
	}
}
