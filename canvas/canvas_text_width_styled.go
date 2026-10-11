package canvas

// TextWidthStyled is how wide a string is under a style in this canvas's face.
func (cv *Canvas) TextWidthStyled(text string, o TextStyle) int {
	return faceWidth(o.face(cv.face()), text, o)
}
