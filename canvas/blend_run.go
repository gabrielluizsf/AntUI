package canvas

// blendRun composites a translucent colour over a run of pixels, and records
// the part of the run that really moved: a blend has no answer to give before
// it happens.
func (cv *Canvas) blendRun(y, x0, x1 int, c Color) {
	// On a layer the run is composited rather than written, since what is under
	// it may be partly transparent; the solid write below would flatten it. A
	// fully opaque colour is the same arithmetic either way, so the fast path
	// still takes it.
	if cv.over && c.A() != 255 {
		row := cv.Pixels[y*cv.Stride+x0 : y*cv.Stride+x1]
		for i := range row {
			row[i] = BlendOver(row[i], c)
		}
		cv.markRun(y, x0, x1)
		return
	}
	BlendRowSolid(cv.Pixels[y*cv.Stride+x0:y*cv.Stride+x1], c)
	cv.markRun(y, x0, x1)
}
