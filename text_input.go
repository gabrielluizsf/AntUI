package antui

// TextInput is the text typed this frame, in UTF-8, empty when nothing was.
func (win *Window) TextInput() string { return win.text.String() }
