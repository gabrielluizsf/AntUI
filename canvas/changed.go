package canvas

// Changed is the rectangle that covered every pixel that really moved when the
// last frame was committed by [Canvas.PresentChanges]: the part of the screen
// that had to be sent for it. It is empty before the first commit, and empty
// again after one that changed nothing.
func (cv *Canvas) Changed() (Area, bool) {
	return cv.change, cv.change.Width > 0
}
