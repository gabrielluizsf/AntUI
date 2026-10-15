package canvas

// Triangle draws the three edges of a triangle.
func (cv *Canvas) Triangle(points [3]Point, c Color) {
	p0, p1, p2 := points[0], points[1], points[2]
	x0, y0 := int(p0.X), int(p0.Y)
	x1, y1 := int(p1.X), int(p1.Y)
	x2, y2 := int(p2.X), int(p2.Y)
	cv.Line(Point{X: float64(x0), Y: float64(y0)}, Point{X: float64(x1), Y: float64(y1)}, c)
	cv.Line(Point{X: float64(x1), Y: float64(y1)}, Point{X: float64(x2), Y: float64(y2)}, c)
	cv.Line(Point{X: float64(x2), Y: float64(y2)}, Point{X: float64(x0), Y: float64(y0)}, c)
}