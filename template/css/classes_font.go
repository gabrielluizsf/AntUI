package css

// Font is the @font-face a font-family list asks for at a weight and a
// slant, or nil when the table's sheet holds none of the names — which is
// when the canvas keeps the face it already draws with.
func (c *CSSClasses) Font(family string, weight uint16, slanted bool) *FontFace {
	if c == nil || c.sheet == nil {
		return nil
	}
	return c.sheet.Font(family, weight, slanted)
}
