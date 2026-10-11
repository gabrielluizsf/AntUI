package canvas

// markRegion is markRun for a whole rectangle, which is what an in-place pass
// — a blur, a filter, a fade — hands over: it rewrote every pixel of it.
func (cv *Canvas) markRegion(r Area) {
	if cv.base == nil {
		return
	}
	x0, y0 := max(r.X, 0), max(r.Y, 0)
	x1, y1 := min(r.X+r.Width, cv.Width), min(r.Y+r.Height, cv.Height)
	for y := y0; y < y1; y++ {
		cv.markRun(y, x0, x1)
	}
}
