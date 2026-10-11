package antui

// WidgetClick drives the press-and-release logic every clickable widget
// shares and reports whether it was clicked this frame. See [Window.Clicked].
func (win *Window) WidgetClick(id uint32, hovered bool) bool {
	return win.widgetClick(id, hovered)
}
