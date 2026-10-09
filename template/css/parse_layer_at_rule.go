package css

// parseLayerAtRule reads @layer and @container: names for where a rule sits
// in a sheet a browser folds. A canvas has neither to fold, so the rules
// inside apply as if the name were the only one. A statement form
// (@layer a, b;) drops it.
func (p *parser) parseLayerAtRule(sh *Sheet, outer Media) error {
	_, semi, err := p.readHeader(true)
	if err != nil {
		return err
	}
	if semi {
		if p.peek() == ';' {
			p.next()
		}
		return nil
	}
	if p.peek() == '{' {
		p.next()
	}
	return p.parseRuleGroup(sh, outer)
}
