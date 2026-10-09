package css

// Keyframes returns the @keyframes block with the given name, the definition
// the stylesheet gave that an animation references, or nil.
func (c *CSSClasses) Keyframes(name string) *Keyframes {
	if c.sheet == nil {
		return nil
	}
	return c.sheet.Keyframes(name)
}
