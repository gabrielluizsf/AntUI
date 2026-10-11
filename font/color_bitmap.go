package font

// colorBitmap is where one glyph's image data begins in the CBDT table, and
// the metrics when they are the index's own: image format 19 keeps none of
// its own and is written against an index that has them.
type colorBitmap struct {
	off, length int
	imageFormat int
	metrics     sbitMetrics
	inIndex     bool
}
