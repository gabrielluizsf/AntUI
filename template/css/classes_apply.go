package css

// Apply tags the table with one or more raw rules (a selector followed by a
// declaration block) folded into the sheet as if they were later than the
// file, so they win the cascade. Parse errors are reported and the rules are
// dropped; media queries inside are honoured.
func (c *CSSClasses) Apply(rules string) error {
	if c.sheet == nil {
		return errNoSheet
	}
	sh, err := Parse(rules)
	if err != nil {
		return err
	}
	if len(sh.rules) == 0 {
		return errEmptyRule
	}
	for _, r := range sh.rules {
		r.Order = c.sheet.order
		c.sheet.order++
	}
	c.sheet.rules = append(c.sheet.rules, sh.rules...)
	c.sheet.Warn = append(c.sheet.Warn, sh.Warn...)
	c.specs = make(map[styleKey]Style)
	return nil
}
