package canvas

// At reads a pixel, in any format. Points outside the canvas read as zero.
// Unlike drawing, reading ignores the clip.
func (cv *Canvas) At(x, y int) Color {
	if x < 0 || y < 0 || x >= cv.Width || y >= cv.Height {
		return 0
	}
	if cv.Pixels == nil {
		if cv.narrow == nil {
			return 0
		}
		return cv.narrowRead(x, y)
	}
	return cv.Pixels[y*cv.Stride+x]
}
