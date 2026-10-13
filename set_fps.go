package antui

import "time"

// SetFPS caps the frame rate. Zero or less — the default — follows the
// display, so a window is paced by the screen it is on; a positive number is a
// ceiling whatever the display says. A window that used to run flat out between
// Begin and End is now paced unless it asks otherwise.
func (win *Window) SetFPS(fps int) {
	win.targetFPS = max(fps, 0)
	// The old grid was set for the old rate; it means nothing now.
	win.nextFrame = time.Time{}
}
