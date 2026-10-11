package antui

// WidgetActive is the widget that is in the middle of being pressed, or 0.
func (win *Window) WidgetActive() uint32 { return win.uiActive }
