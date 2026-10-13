package canvas

// SetClip narrows drawing to a rectangle, intersected with the canvas.
func (cv *Canvas) SetClip(a Area) {
	x, y, w, h := a.X, a.Y, a.Width, a.Height
	x0 := clampInt(x, 0, cv.Width)
	y0 := clampInt(y, 0, cv.Height)
	x1 := clampInt(x+w, 0, cv.Width)
	y1 := clampInt(y+h, 0, cv.Height)
	cv.Clip = Area{x0, y0, max(0, x1-x0), max(0, y1-y0)}
}
