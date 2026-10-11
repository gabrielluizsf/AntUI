package font

// simpleGlyph reads a glyph whose outline is its own points: the raw table
// bytes, where each contour ends, and the per-point flags once they are read.
type simpleGlyph struct {
	data   []byte
	ends   []int
	flags  []byte
	points int
}
