package css

// parseRuleGroup reads the body of a rule-group at-rule — @media, @supports,
// @layer, @container — which contains whole rules and nested at-rules (a
// @media inside a @media, a rule directly inside @supports), until the
// closing '}'.
func (p *parser) parseRuleGroup(sh *Sheet, m Media) error {
	for {
		p.skipSpace()
		if p.eof() {
			return fmtErrf("unterminated block, missing '}'")
		}
		switch {
		case p.peek() == '}':
			p.next()
			return nil
		case p.peek() == ';':
			p.next()
		case p.peek() == '@':
			if err := p.parseAtRule(sh, m); err != nil {
				return err
			}
		default:
			if err := p.parseRule(sh, m); err != nil {
				return err
			}
		}
	}
}
