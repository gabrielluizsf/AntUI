package canvas

// TextScaled draws a string larger, and returns the width of the widest line.
func (cv *Canvas) TextScaled(pos Point, text string, o TextStyle) int {
	scale := o.Scale
	if scale < 1 {
		scale = 1
	}
	c := o.Color
	f := cv.face()
	if f != BuiltinFace() {
		f = f.Scaled(scale)
		return f.Draw(FaceDrawParams{Canvas: cv, Pos: Point{X: pos.X, Y: float64(int(pos.Y) + f.Ascent())}, Text: text, Color: c})
	}
	return cv.textBitmap(int(pos.X), int(pos.Y), text, c, float64(scale))
}
