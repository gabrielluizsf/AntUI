package canvas

// Text draws a string and returns the width of the widest line drawn.
func (cv *Canvas) Text(pos Point, text string, c Color) int {
	return cv.TextScaled(pos, text, TextStyle{Color: c})
}
