package canvas

// Triangle draws the three edges of a triangle.
func (cv *Canvas) Triangle(x0, y0, x1, y1, x2, y2 int, c Color) {
	cv.Line(x0, y0, x1, y1, c)
	cv.Line(x1, y1, x2, y2, c)
	cv.Line(x2, y2, x0, y0, c)
}
