package canvas

// filterBounds intersects a rectangle with the canvas and its active clip.
func (cv *Canvas) filterBounds(x, y, w, h int) (x0, y0, x1, y1 int) {
	x0 = max(x, cv.Clip.X)
	y0 = max(y, cv.Clip.Y)
	x1 = min(x+w, cv.Clip.X+cv.Clip.Width)
	y1 = min(y+h, cv.Clip.Y+cv.Clip.Height)
	x0 = max(x0, 0)
	y0 = max(y0, 0)
	x1 = min(x1, cv.Width)
	y1 = min(y1, cv.Height)
	if x0 >= x1 || y0 >= y1 {
		return 0, 0, 0, 0
	}
	return x0, y0, x1, y1
}
