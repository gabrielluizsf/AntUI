package canvas

// StrokeStyle is everything about a stroke that is not its colour: how wide, how
// the ends and the corners are treated, and how far past that limit a miter may
// grow before it is beveled instead.
type StrokeStyle struct {
	Width      float64
	Cap        LineCap
	Join       LineJoin
	MiterLimit float64
}
