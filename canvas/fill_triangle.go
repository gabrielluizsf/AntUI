package canvas

// FillTriangle fills a triangle. Either winding works.
func (cv *Canvas) FillTriangle(x0, y0, x1, y1, x2, y2 int, c Color) {
	minX := max(min(x0, min(x1, x2)), cv.Clip.X)
	maxX := min(max(x0, max(x1, x2)), cv.Clip.X+cv.Clip.Width-1)
	minY := max(min(y0, min(y1, y2)), cv.Clip.Y)
	maxY := min(max(y0, max(y1, y2)), cv.Clip.Y+cv.Clip.Height-1)

	if area := (x1-x0)*(y2-y0) - (x2-x0)*(y1-y0); area == 0 {
		return
	}

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			w0 := (x1-x0)*(y-y0) - (x-x0)*(y1-y0)
			w1 := (x2-x1)*(y-y1) - (x-x1)*(y2-y1)
			w2 := (x0-x2)*(y-y2) - (x-x2)*(y0-y2)
			if (w0 >= 0 && w1 >= 0 && w2 >= 0) || (w0 <= 0 && w1 <= 0 && w2 <= 0) {
				cv.Pixel(x, y, c)
			}
		}
	}
}
