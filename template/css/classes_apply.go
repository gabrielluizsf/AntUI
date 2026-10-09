package css

// Apply tags the table with a raw rule (a selector followed by a declaration
// block) folded into the sheet as if it were another rule later than the
// file, so it wins the cascade. Parse errors are reported and the rule is
// dropped; media queries inside are honoured.
func (c *CSSClasses) Apply(rule string) error {
	if c.sheet == nil {
		return errNoSheet
	}
	m := &parser{src: rule}
	sh := &Sheet{order: c.sheet.order}
	if err := m.parseRule(sh, Media{}); err != nil {
		return err
	}
	if len(sh.rules) == 0 {
		return errEmptyRule
	}
	c.sheet.rules = append(c.sheet.rules, sh.rules...)
	c.sheet.order = sh.order
	c.sheet.Warn = append(c.sheet.Warn, sh.Warn...)
	c.specs = make(map[styleKey]Style)
	return nil
}
