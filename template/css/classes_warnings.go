package css

// Warnings are the sheet's own: every declaration the engine did not
// understand, for the template to show its author.
func (c *CSSClasses) Warnings() []string {
	if c.sheet == nil {
		return nil
	}
	return c.sheet.Warn
}
