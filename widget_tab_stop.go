package antui

// WidgetTabStop says the widget takes the keyboard, at its place in the
// frame's focus order. Widgets register themselves, in the order they draw,
// and Tab walks that order.
func (win *Window) WidgetTabStop(id uint32) { win.uiTab = append(win.uiTab, id) }
