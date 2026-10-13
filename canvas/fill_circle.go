package canvas

// FillCircle draws a filled, antialiased disc.
func (cv *Canvas) FillCircle(c Circle, col Color) {
	cv.FillEllipse(Ellipse{X: c.X, Y: c.Y, RX: c.R, RY: c.R}, col)
}
