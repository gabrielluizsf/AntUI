package canvas

// Rect draws a one-pixel outline.
func (cv *Canvas) Rect(x, y, w, h int, c Color) {
	if w <= 0 || h <= 0 {
		return
	}
	cv.FillRect(x, y, w, 1, c)
	if h > 1 {
		cv.FillRect(x, y+h-1, w, 1, c)
	}
	if h > 2 {
		cv.FillRect(x, y+1, 1, h-2, c)
		if w > 1 {
			cv.FillRect(x+w-1, y+1, 1, h-2, c)
		}
	}
}
