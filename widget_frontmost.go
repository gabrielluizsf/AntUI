package antui

// WidgetFrontmost defers draw until the very end of the frame, after every
// widget painted in flow order, so an open dropdown or calendar paints above
// whatever sits under it. The window runs the draws just before presenting,
// and forgets them at the next Begin, so a caller registers one every frame.
func (win *Window) WidgetFrontmost(draw func()) {
	win.frontmost = append(win.frontmost, draw)
}
