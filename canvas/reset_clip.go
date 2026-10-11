package canvas

// ResetClip opens drawing back up to the whole canvas.
func (cv *Canvas) ResetClip() {
	cv.Clip = Area{0, 0, cv.Width, cv.Height}
}
