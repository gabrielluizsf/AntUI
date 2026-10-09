package css

// parseMediaAtRule reads @media's prelude, folds it into the enclosing media
// and parses the rules inside under the combined query. A @media inside a
// @media must also live inside the outer one.
func (p *parser) parseMediaAtRule(sh *Sheet, outer Media) error {
	prelude, _, err := p.readHeader(false)
	if err != nil {
		return err
	}
	if p.peek() == '{' {
		p.next()
	}
	inner, warns := parseMedia(prelude)
	sh.Warn = append(sh.Warn, warns...)
	return p.parseRuleGroup(sh, outer.and(inner))
}
