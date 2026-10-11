package antui

// resolveGeometry applies the display scale and size limits to the freshly
// opened window.
func (win *Window) resolveGeometry(opt Options) {
	if scale, ok := win.native.displayScale(); ok {
		win.scale = scale
	}
	width, height := opt.Width, opt.Height
	if opt.Points && win.scale > 0 {
		width, height = scaleTo(width, win.scale), scaleTo(height, win.scale)
		win.limits = win.limits.Scaled(win.scale)
	}
	if dw, dh, ok := win.native.displaySize(); ok {
		width, height = fitDisplay(width, height, dw, dh)
	}
	width, height = win.limits.Clamp(width, height)

	if width != opt.Width || height != opt.Height {
		win.native.setSize(win, width, height)
		win.ResizeCanvas(width, height)
		win.windowedWidth, win.windowedHeight = width, height
	}
	win.native.setLimits(win, win.Bounds())
}
