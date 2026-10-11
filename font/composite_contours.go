package font

// compositeContours reads a glyph made of other glyphs — an accented letter
// is a letter and an accent, placed.
func (f *TTF) compositeContours(data []byte, depth int) []Contour {
	var out []Contour
	p := 10
	for {
		comp, next, ok := readComponent(data, p)
		if !ok {
			return out
		}
		p = next
		for _, part := range f.GlyphContours(comp.index, depth+1) {
			out = append(out, transformContour(comp, part))
		}
		if !comp.more {
			return out
		}
	}
}
