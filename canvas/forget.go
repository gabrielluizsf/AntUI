package canvas

// Forget drops the baseline, so writes go straight to the pixels again with
// nothing recorded. The Dirty bounds still stand: they say what was written,
// which is a different question from what changed.
func (cv *Canvas) Forget() {
	cv.base = nil
	cv.rows = cv.rows[:0]
	cv.change = Area{}
}
