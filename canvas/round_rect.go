package canvas

// RoundRect draws a one-pixel outline with antialiased corners, each a quarter
// ellipse with the radius [Rounded] gives across and down. Equal radii — what
// almost every rounded rectangle in a UI is — come out as circles, and the
// radii are pulled in to half the box in each direction.
func (cv *Canvas) RoundRect(b Rounded, c Color) {
	x, y, w, h := b.X, b.Y, b.Width, b.Height
	rx, ry := min(b.RX, w/2), min(b.vertical(), h/2)
	if w <= 0 || h <= 0 {
		return
	}
	if rx <= 0 || ry <= 0 {
		cv.Rect(b.Area, c)
		return
	}

	cv.FillRect(Area{x + rx, y, w - 2*rx, 1}, c)
	cv.FillRect(Area{x + rx, y + h - 1, w - 2*rx, 1}, c)
	cv.FillRect(Area{x, y + ry, 1, h - 2*ry}, c)
	cv.FillRect(Area{x + w - 1, y + ry, 1, h - 2*ry}, c)

	roundRectCorners{
		cv: cv,
		c:  c,
		rx: rx,
		ry: ry,
		corners: [4][2]int{
			{x + rx - 1, y + ry - 1},
			{x + w - rx, y + ry - 1},
			{x + rx - 1, y + h - ry},
			{x + w - rx, y + h - ry},
		},
	}.draw()
}
