package font

// Outline is where one glyph's contours are in the glyf table.
func (f *TTF) Outline(glyph int) []byte {
	if glyph < 0 || glyph >= f.NumGlyphs {
		return nil
	}
	var start, end int
	if f.LongLoca {
		if (glyph+2)*4 > len(f.Loca) {
			return nil
		}
		start, end = int(be32(f.Loca, glyph*4)), int(be32(f.Loca, (glyph+1)*4))
	} else {
		if (glyph+2)*2 > len(f.Loca) {
			return nil
		}
		start, end = int(be16(f.Loca, glyph*2))*2, int(be16(f.Loca, (glyph+1)*2))*2
	}
	if start >= end || end > len(f.Glyf) {
		return nil // an empty glyph, which a space is
	}
	return f.Glyf[start:end]
}
