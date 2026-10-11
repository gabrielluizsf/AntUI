package canvas

// span is a run of pixels in one row written at one coverage, which is all the
// fill has to say about a pixel: how much of it the shape covers.
type span struct {
	x0, x1 int
	y      int
	a      float64
}
