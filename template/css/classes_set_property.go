package css

// SetProperty overrides a property for every widget of the given tag and
// classes, appended after the sheet — so it wins whatever the file said for
// those selectors — and clears the cache.
func (c *CSSClasses) SetProperty(q PropertyQuery, raw string) {
	if c.sheet == nil {
		return
	}
	c.sheet.SetProperty(q, raw)
	c.specs = make(map[styleKey]Style)
}
