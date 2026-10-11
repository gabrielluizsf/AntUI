package canvas

// DrawStyled writes a string with a [TextStyle], the baseline of the first
// line at y, and returns the width of the widest line. It is [Face.Draw] with
// the weight, slant and spacing the stylesheet asked for.
func (cv *Canvas) DrawStyled(x, y int, text string, c Color, o TextStyle) int {
	o = o.norm()
	f := o.face(cv.face())
	if f != BuiltinFace() && f.file != nil {
		sf := f.Scaled(o.Scale)
		return cv.drawStyledMask(sf, x, y+sf.Ascent(), text, c, o)
	}
	return cv.drawStyledBitmap(x, y, text, c, o)
}
