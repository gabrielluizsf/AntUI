package canvas

// markRun is markSpan for a path that wrote a run of pixels between x0 and
// x1 without knowing what each of them became — a blend, a copy, a filter.
// The run is half-open, as the write loops that call it are.
func (cv *Canvas) markRun(y, x0, x1 int) {
	if x1 > x0 {
		cv.markSpan(y, x0, x1-1)
	}
}
