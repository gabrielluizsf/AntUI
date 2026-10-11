package antui

// WidgetKeyActivate reports whether Enter or Space was pressed this frame
// while the widget held the keyboard. A focused button, checkbox or radio is
// "clicked" by Enter or Space, the way the web does it.
func (win *Window) WidgetKeyActivate(id uint32) bool {
	if win.uiFocus != id || (win.uiActive != 0 && win.uiActive != id) {
		return false
	}
	return win.KeyPressed(KeyEnter) || win.KeyPressed(KeySpace)
}
