package canvas

// Comparing reports whether this canvas is measuring its writes against a
// baseline.
func (cv *Canvas) Comparing() bool { return cv.base != nil }
