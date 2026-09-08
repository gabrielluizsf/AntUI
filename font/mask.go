package font

// Mask is one glyph rendered to coverage, 0 to 255 per pixel.
type Mask struct {
	W, H   int
	Left   int
	Top    int
	Alpha  []uint8
	Stride int
}
