package css

// SetStyle loads a CSS file into the table, replacing whatever it held, and
// clears the style cache. Compile errors — a malformed selector or an
// unbalanced block — stop the load; properties the engine does not know are
// skipped and reported through the sheet's warnings.
func (c *CSSClasses) SetStyle(cssFile string) error {
	if c.specs == nil {
		c.specs = make(map[styleKey]Style)
	}
	sh, err := ParseFile(cssFile)
	if err != nil {
		return err
	}
	c.sheet = sh
	c.specs = make(map[styleKey]Style)
	return nil
}
