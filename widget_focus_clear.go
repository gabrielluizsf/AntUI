package antui

// WidgetFocusClear takes the keyboard away again.
func (win *Window) WidgetFocusClear(id uint32) {
	if win.uiFocus == id {
		win.uiFocus = 0
	}
}
