package font

// Contour is one closed loop of a glyph, in font units, with each point
// marked as on or off the curve.
type Contour struct {
	X, Y []float64
	On   []bool
}
