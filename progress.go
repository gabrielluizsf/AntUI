package antui

import "github.com/gabrielluizsf/antui/canvas"

// Progress draws a progress bar, with progress running from 0 to 1.
func (win *Window) Progress(area canvas.Area, progress float32) {
	x, y, w, h := area.X, area.Y, area.Width, area.Height
	if w <= 0 || h <= 0 {
		return
	}
	progress = min(max(progress, 0), 1)
	filled := int(float32(w)*progress + 0.5)

	win.cv.FillRoundRect(canvas.Round(area, h/2), win.theme.Border)
	if filled > 0 {
		win.cv.FillRoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: filled, Height: h}, h/2), win.theme.Accent)
	}
}
