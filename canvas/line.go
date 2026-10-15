package canvas

// Line draws a line with Bresenham's algorithm.
func (cv *Canvas) Line(a, b Point, c Color) {
	x0, y0 := int(a.X), int(a.Y)
	x1, y1 := int(b.X), int(b.Y)
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := -1, -1
	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy
	for {
		cv.Pixel(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}
