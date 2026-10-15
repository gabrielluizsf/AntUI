package canvas

// fillCoverageSpan paints one run. A run at the colour's own alpha goes through
// the fastest write the canvas has, a partial one scales the alpha and goes a
// pixel at a time, because a blend is a blend of one pixel and not of a run.
func (cv *Canvas) fillCoverageSpan(s span, c Color) {
	a := uint8(float64(c.A())*s.a + 0.5)
	if a == 0 {
		return
	}
	switch {
	case a != c.A():
		sc := RGBA(RGBAParams{R: c.R(), G: c.G(), B: c.B(), A: a})
		for x := s.x0; x < s.x1; x++ {
			cv.Pixel(x, s.y, sc)
		}
	case cv.Pixels != nil:
		cv.blendRun(s.y, s.x0, s.x1, c)
	default:
		cv.narrowRun(s.x0, s.x1, s.y, c)
		cv.markRun(s.y, s.x0, s.x1)
	}
}
