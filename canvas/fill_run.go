package canvas

// fillRun writes an opaque colour over a run of pixels and records the part of
// the run that moved. Deciding not to write a run that already holds the colour
// would cost about as much as writing it — a comparison is a load and a compare
// per pixel, and a store is a load and a store — so the run is written and then
// measured against the frame on the screen, which does reject a whole frame in
// one memcmp.
func (cv *Canvas) fillRun(y, x0, x1 int, c Color) {
	solidRunGo(cv.Pixels[y*cv.Stride+x0:y*cv.Stride+x1], c)
	cv.markRun(y, x0, x1)
}
