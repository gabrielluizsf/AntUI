package font

// strikeReader accumulates the glyphs of one IndexSubTable. constant is the
// one set of metrics every glyph of the range shares, kept here rather than
// repeated in each glyph's data.
type strikeReader struct {
	cblc, cbdt  []byte
	sub         int
	first, last int
	imageFormat int
	base        int
	out         map[int]colorBitmap
	constant    sbitMetrics
	hasConst    bool
}
