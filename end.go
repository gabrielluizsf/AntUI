package antui

import ()

// End finishes the frame: it sends the pixels that changed to the screen, then
// holds the loop to the frame rate.
//
// It reports whether anything reached the display. A frame that redrew the
// same picture — a window sitting still with no animation running — sends
// nothing and answers false, and so costs almost nothing: no upload, no scan
// of the canvas, no drawing the backend had to be woken for.
func (win *Window) End() bool {
	win.moveFocus()
	// Drawings registered while the frame was being painted land here, above
	// everything, before the pixels are presented.
	for i := len(win.frontmost) - 1; i >= 0; i-- {
		win.frontmost[i]()
	}
	if win.native == nil {
		return false
	}
	dirty := win.dirtyRegion()
	// A frame that drew the picture already on the screen is not presented at
	// all: the platform is not asked, which is what keeps an idle window from
	// waking the compositor sixty times a second to be told there is nothing to
	// do. The canvas's own copy of the screen has already been brought up to
	// date by the measuring, so the next frame starts from the right place.
	if dirty.Width > 0 && dirty.Height > 0 {
		win.native.present(win, dirty)
		win.presented = true
	} else {
		win.presented = false
	}
	win.pace()
	return win.presented
}
