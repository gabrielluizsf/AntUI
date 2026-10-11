package canvas

// TextWidthStyled is [TextWidth] under a [TextStyle], measured with the
// program's default face.
func TextWidthStyled(text string, o TextStyle) int {
	return faceWidth(o.face(DefaultFace()), text, o)
}
