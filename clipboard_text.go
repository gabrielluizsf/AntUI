package antui

// ClipboardText is what the system clipboard holds as text, or "" when it
// holds nothing this window can read.
func (win *Window) ClipboardText() string {
	if win == nil || win.native == nil {
		return ""
	}
	text, _ := win.native.clipboard()
	return text
}
