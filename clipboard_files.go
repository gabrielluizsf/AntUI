package antui

// ClipboardFiles is the files the clipboard names. Empty when it holds none.
func (win *Window) ClipboardFiles() []string {
	if win == nil || win.native == nil {
		return nil
	}
	_, files := win.native.clipboard()
	return files
}
