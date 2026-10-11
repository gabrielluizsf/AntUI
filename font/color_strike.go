package font

// colorStrike is one size the file keeps bitmaps at.
type colorStrike struct {
	ppem   int
	glyphs map[int]colorBitmap
}
