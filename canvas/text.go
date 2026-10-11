package canvas

// Text draws a string and returns the width of the widest line drawn.
func (cv *Canvas) Text(x, y int, text string, c Color) int {
	return cv.TextScaled(x, y, text, c, 1)
}
