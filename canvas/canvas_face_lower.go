package canvas

// face is the face to use, never nil.
func (cv *Canvas) face() *Face {
	if cv != nil && cv.hasFont && cv.font != nil {
		return cv.font
	}
	return DefaultFace()
}
