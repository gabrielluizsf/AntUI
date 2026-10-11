package antui

// fitDisplay reduces a size to what the display can actually show.
//
// Asking for a window larger than the screen is not an error and must not be
// one — a 1920 by 1080 game opened on a 1366 by 768 laptop is an ordinary
// Tuesday. The request is reduced to what fits and the game is told what it
// got, by Width and Height reporting the truth.
func fitDisplay(width, height, displayWidth, displayHeight int) (int, int) {
	if displayWidth > 0 && width > displayWidth {
		width = displayWidth
	}
	if displayHeight > 0 && height > displayHeight {
		height = displayHeight
	}
	return width, height
}
