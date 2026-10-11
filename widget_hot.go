package antui

// WidgetHot marks the widget as the one under the pointer, which is how a
// widget that is not drawn under the cursor still keeps others from acting on
// presses meant for it.
func (win *Window) WidgetHot(id uint32) { win.uiHot = id }
