package canvas

// SetFace makes this one canvas draw text in a face of its own, whatever the
// program's default is. nil goes back to the default.
func (cv *Canvas) SetFace(f *Face) {
	if cv != nil {
		cv.font = f
		cv.hasFont = true
		if f == nil {
			cv.hasFont = false
		}
	}
}
