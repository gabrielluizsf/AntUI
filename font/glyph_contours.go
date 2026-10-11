package font

// GlyphContours reads one glyph's outline, following composites.
func (f *TTF) GlyphContours(glyph, depth int) []Contour {
	if depth > 5 {
		return nil
	}
	data := f.Outline(glyph)
	if len(data) < 10 {
		return nil
	}
	if n := int16(be16(data, 0)); n >= 0 {
		return simpleContours(data, int(n))
	}
	return f.compositeContours(data, depth)
}
