package font

// TTF is a parsed font file. One file may hold several faces at several
// sizes; this is the file, and [Face] is one size of it.
type TTF struct {
	Data []byte

	UnitsPerEm uint16
	LongLoca   bool
	NumGlyphs  int

	Ascent, Descent, LineGap int16
	NumHMetrics              int

	Loca, Glyf, Hmtx []byte
	Cmap             CmapTable

	// cbdt holds the colour bitmaps and colorStrikes the CBLC index reads
	// into: a picture per glyph instead of an outline, which is what an
	// emoji font is made of. Both are empty when the file carries neither.
	cbdt         []byte
	colorStrikes []colorStrike
}
