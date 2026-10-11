package font

// Glyph is the glyph index for a codepoint, or 0 — which is the font's own
// "not there" glyph, usually an empty box, and is what should be drawn.
func (t CmapTable) Glyph(r rune) int {
	switch t.Format {
	case 4:
		return t.glyph4(r)
	case 12:
		return t.glyph12(r)
	case 6:
		return t.glyph6(r)
	case 0:
		return t.glyph0(r)
	}
	return 0
}
