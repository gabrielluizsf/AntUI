package antui

// WidgetCaret reports whether the text cursor is lit right now. It blinks on
// the same clock the built-in text field blinks on.
func (win *Window) WidgetCaret() bool { return int64(win.uiBlink*2)%2 == 0 }
