package canvas

// FillRoundRect fills a rectangle with antialiased corners, each a quarter
// ellipse with the radius [Rounded] gives across and down. Equal radii — what
// almost every rounded rectangle in a UI is — come out as circles, and the
// radii are pulled in to half the box in each direction so a radius wider than
// the box does not turn it inside out.
func (cv *Canvas) FillRoundRect(b Rounded, c Color) {
	x, y, w, h := b.X, b.Y, b.Width, b.Height
	rx, ry := min(b.RX, w/2), min(b.vertical(), h/2)
	if w <= 0 || h <= 0 {
		return
	}
	if rx <= 0 || ry <= 0 {
		cv.FillRect(b.Area, c)
		return
	}
	baseAlpha := int(c.A())
	if baseAlpha == 0 {
		return
	}

	cv.FillRect(Area{x + rx, y, w - 2*rx, h}, c)
	cv.FillRect(Area{x, y + ry, rx, h - 2*ry}, c)
	cv.FillRect(Area{x + w - rx, y + ry, rx, h - 2*ry}, c)

	fillRoundRectCorners{
		cv: cv, c: c, x: x, y: y, w: w, h: h, rx: rx, ry: ry, alpha: baseAlpha,
	}.draw()
}
