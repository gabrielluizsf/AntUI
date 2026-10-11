package antui

// WidgetEdit applies this frame's keys and typing to the focused field,
// exactly as the built-in Input does.
func (win *Window) WidgetEdit(text *string) bool { return win.editText(text) }
