package css

// Property reads the winning raw value of one property for a tag and class
// list, computed the way the cascade would compute it. It answers what the
// stylesheet says, outside of the folded [Style] the drawing uses, reading
// media queries against a window as tall as width is wide — the drawing goes
// through [CSSClasses.GetStyleViewport], which measures a real one.
func (c *CSSClasses) Property(tag string, classes []string, prop string, width int) (string, bool) {
	if c.sheet == nil {
		return "", false
	}
	return c.sheet.Property(tag, classes, prop, width)
}
