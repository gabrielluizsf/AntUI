package antui

// SetExternal hands the window's surface to another renderer, or takes it
// back. While on, [Window.End] stops copying the canvas into the window — it
// still pumps input and paces the loop — and the platform stops painting its
// own background over the other renderer's frames. Turning it off sends the
// whole canvas again on the next frame, so falling back to software drawing
// takes over cleanly from whatever the other renderer left there.
//
// It reports whether the platform could hand the surface over. On a platform
// that cannot, the canvas is still kept off the window, but the platform may
// repaint over the other renderer.
func (win *Window) SetExternal(on bool) bool {
	if win.external && !on {
		win.fullRedraw = true
	}
	win.external = on
	if win.native == nil {
		return false
	}
	return win.native.setExternal(on)
}
