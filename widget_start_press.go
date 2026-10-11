package antui

// WidgetStartPress begins a press on the widget. It takes the keyboard away
// from any text field, the same way the built-in widgets do when pressed.
func (win *Window) WidgetStartPress(id uint32) {
	win.uiActive = id
	win.uiFocus = 0
}
