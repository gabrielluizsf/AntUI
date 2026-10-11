package antui

// WidgetFocusSet hands the keyboard to the widget.
func (win *Window) WidgetFocusSet(id uint32) {
	win.uiFocus = id
	win.uiBlink = 0
}
