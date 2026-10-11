//go:build android

package antui

// sync makes the canvas the size of the surface. On a desktop a resize is
// news from the window manager; here it is simply a fact to be noticed,
// because nothing was ever asked for.
func (b *androidWindow) sync(win *Window) {
	if b.surface == nil {
		return
	}
	w, h := b.surface.Width(), b.surface.Height()
	if w <= 0 || h <= 0 || (w == win.cv.Width && h == win.cv.Height) {
		return
	}
	if win.ResizeCanvas(w, h) {
		win.Push(Event{Type: EventResize, Width: w, Height: h})
		// resizeCanvas throws the safe area away, since it was in the old
		// canvas's coordinates. The platform's rectangle is still current.
		b.safe(win)
	}
}
