package canvas

// Clear paints the whole clip region. The colour is forced opaque.
func (cv *Canvas) Clear(c Color) {
	cv.FillRect(cv.Clip.X, cv.Clip.Y, cv.Clip.Width, cv.Clip.Height, c|0xFF000000)
}
