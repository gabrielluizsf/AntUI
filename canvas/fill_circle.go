package canvas

// FillCircle draws a filled, antialiased disc.
func (cv *Canvas) FillCircle(cx, cy, radius int, c Color) {
	cv.FillEllipse(cx, cy, radius, radius, c)
}
