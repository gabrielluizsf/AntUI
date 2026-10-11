package canvas

// markSpan remembers that a row was written between x0 and x1. The span is
// the widest the row was written, which is all the measuring needs: the run is
// narrowed to the pixels that really moved when the frame is committed.
func (cv *Canvas) markSpan(y, x0, x1 int) {
	if cv.base == nil || x1 < x0 || y < 0 || y >= len(cv.spans) {
		return
	}
	s := &cv.spans[y]
	if s.left < 0 {
		s.left, s.right = x0, x1
		cv.rows = append(cv.rows, y)
		return
	}
	s.left, s.right = min(s.left, x0), max(s.right, x1)
}
