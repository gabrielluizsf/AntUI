package antui

import "math"

// SizeInPoints is the size the window appears to be, as against the pixels it
// draws — the two are the same number until the display is scaled.
func (win *Window) SizeInPoints() (width, height int) {
	scale := win.ContentScale()
	if win == nil {
		return 0, 0
	}
	return max(int(math.Round(float64(win.Width())/scale)), 1),
		max(int(math.Round(float64(win.Height())/scale)), 1)
}
