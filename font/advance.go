package font

// Advance is how far the pen moves after a glyph, in font units.
func (f *TTF) Advance(glyph int) int16 {
	if f.NumHMetrics == 0 || len(f.Hmtx) < 4 {
		return int16(f.UnitsPerEm / 2)
	}
	i := min(glyph, f.NumHMetrics-1)
	if (i+1)*4 > len(f.Hmtx) {
		return int16(f.UnitsPerEm / 2)
	}
	return int16(be16(f.Hmtx, i*4))
}
