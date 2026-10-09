package css

// parseSupportsAtRule evaluates @supports' condition against the engine: the
// properties it has, the values it reads, the selectors it can evaluate. A
// block whose test fails is not part of the sheet — that is the test doing
// its job — and one whose condition cannot be read is reported and dropped
// the same way.
func (p *parser) parseSupportsAtRule(sh *Sheet, outer Media) error {
	prelude, semi, err := p.readHeader(true)
	if err != nil {
		return err
	}
	if semi {
		// A @supports with no block gates nothing, and there is nothing
		// inside it to keep.
		if p.peek() == ';' {
			p.next()
		}
		return nil
	}
	ok, warns := supportsCondition(prelude)
	sh.Warn = append(sh.Warn, warns...)
	if !ok {
		return p.skipBlock()
	}
	if p.peek() == '{' {
		p.next()
	}
	return p.parseRuleGroup(sh, outer)
}
