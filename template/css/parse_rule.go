package css

// parseRule reads one selector block: `.a:hover, button { … }`, optionally
// constrained by an enclosing media query.
func (p *parser) parseRule(sh *Sheet, m Media) error {
	selText, semi, err := p.readHeader(true)
	if err != nil {
		return err
	}
	if p.peek() == '{' {
		p.next()
	}
	if semi {
		// A stray selector with no block: dropped, like a browser drops it.
		return nil
	}
	selectors, warns := parseSelectors(selText)
	sh.Warn = append(sh.Warn, warns...)
	rule := &Rule{Media: m, Selectors: selectors, Order: sh.order}
	sh.order++
	if err := p.declBlock(func(text string) error {
		if d, ok := parseDeclaration(text); ok {
			rule.Decls = append(rule.Decls, d)
		}
		return nil
	}); err != nil {
		return err
	}
	sh.rules = append(sh.rules, rule)
	return nil
}
