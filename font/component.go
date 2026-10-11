package font

// component is one placed piece of a composite glyph: the glyph to draw, the
// offset to place it at, the 2x2 transform to draw it under, and whether a
// further component follows.
type component struct {
	index      int
	dx, dy     float64
	a, b, c, d float64
	more       bool
}
