package canvas

// Circle draws a one-pixel circle outline.
func (cv *Canvas) Circle(cx, cy, radius int, c Color) {
	if radius <= 0 {
		return
	}
	x, y := radius, 0
	err := 1 - radius
	for x >= y {
		cv.Pixel(cx+x, cy+y, c)
		cv.Pixel(cx+y, cy+x, c)
		cv.Pixel(cx-y, cy+x, c)
		cv.Pixel(cx-x, cy+y, c)
		cv.Pixel(cx-x, cy-y, c)
		cv.Pixel(cx-y, cy-x, c)
		cv.Pixel(cx+y, cy-x, c)
		cv.Pixel(cx+x, cy-y, c)
		y++
		if err < 0 {
			err += 2*y + 1
		} else {
			x--
			err += 2*(y-x) + 1
		}
	}
}
