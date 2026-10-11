package antui

// WidgetEndPress ends the press on the widget, whatever it was.
func (win *Window) WidgetEndPress(id uint32) {
	if win.uiActive == id {
		win.uiActive = 0
	}
}
