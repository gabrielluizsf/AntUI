package antui

// Progress draws a progress bar, with progress running from 0 to 1.
func (win *Window) Progress(x, y, w, h int, progress float32) {
	if w <= 0 || h <= 0 {
		return
	}
	progress = min(max(progress, 0), 1)
	filled := int(float32(w)*progress + 0.5)

	win.cv.FillRoundRect(x, y, w, h, h/2, win.theme.Border)
	if filled > 0 {
		win.cv.FillRoundRect(x, y, filled, h, h/2, win.theme.Accent)
	}
}
