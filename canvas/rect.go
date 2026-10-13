package canvas

// Rect draws a one-pixel outline.
func (cv *Canvas) Rect(a Area, c Color) {
	x, y, w, h := a.X, a.Y, a.Width, a.Height
	if w <= 0 || h <= 0 {
		return
	}
	cv.FillRect(Area{x, y, w, 1}, c)
	if h > 1 {
		cv.FillRect(Area{x, y + h - 1, w, 1}, c)
	}
	if h > 2 {
		cv.FillRect(Area{x, y + 1, 1, h - 2}, c)
		if w > 1 {
			cv.FillRect(Area{x + w - 1, y + 1, 1, h - 2}, c)
		}
	}
}
