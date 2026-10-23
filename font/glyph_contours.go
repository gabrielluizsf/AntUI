package font

// GlyphContours reads one glyph's outline, following composites.
func (f *TTF) GlyphContours(glyph, depth int) []Contour {
	budget := glyphBudget{components: maxComponents, points: maxGlyphPoints}
	return f.glyphContours(glyph, depth, &budget)
}

// glyphContours is GlyphContours with what the whole glyph may still read,
// which every piece inside it spends from.
func (f *TTF) glyphContours(glyph, depth int, budget *glyphBudget) []Contour {
	if depth > 5 {
		return nil
	}
	data := f.Outline(glyph)
	if len(data) < 10 {
		return nil
	}
	if n := int16(be16(data, 0)); n >= 0 {
		contours := simpleContours(data, int(n))
		if !budget.spendPoints(contours) {
			return nil
		}
		return contours
	}
	return f.compositeContours(data, depth, budget)
}
