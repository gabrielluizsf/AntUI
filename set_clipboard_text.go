package antui

// SetClipboardText puts text on the system clipboard — what a Copy button
// does — and reports whether the system took it.
func (win *Window) SetClipboardText(text string) bool {
	if win == nil || win.native == nil {
		return false
	}
	return win.native.setClipboard(text)
}
