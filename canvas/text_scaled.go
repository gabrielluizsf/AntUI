package canvas

// TextScaled draws a string larger, and returns the width of the widest line.
func (cv *Canvas) TextScaled(x, y int, text string, c Color, scale int) int {
	scale = max(scale, 1)
	f := cv.face()
	if f != BuiltinFace() {
		f = f.Scaled(scale)
		return f.Draw(cv, x, y+f.Ascent(), text, c)
	}
	return cv.textBitmap(x, y, text, c, float64(scale))
}
