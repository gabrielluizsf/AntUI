package antui

// editText applies this frame's keys and typing to the focused field.
func (win *Window) editText(text *string) bool {
	win.uiCursor = min(max(win.uiCursor, 0), len(*text))
	changed := false

	if win.KeyPressed(KeyBackspace) && win.uiCursor > 0 {
		start := utf8Prev(*text, win.uiCursor)
		*text = (*text)[:start] + (*text)[win.uiCursor:]
		win.uiCursor = start
		changed = true
		win.uiBlink = 0
	}
	if win.KeyPressed(KeyDelete) && win.uiCursor < len(*text) {
		end := utf8Next(*text, win.uiCursor)
		*text = (*text)[:win.uiCursor] + (*text)[end:]
		changed = true
		win.uiBlink = 0
	}
	if win.KeyPressed(KeyLeft) {
		win.uiCursor = utf8Prev(*text, win.uiCursor)
		win.uiBlink = 0
	}
	if win.KeyPressed(KeyRight) {
		win.uiCursor = utf8Next(*text, win.uiCursor)
		win.uiBlink = 0
	}
	if win.KeyPressed(KeyHome) {
		win.uiCursor = 0
	}
	if win.KeyPressed(KeyEnd) {
		win.uiCursor = len(*text)
	}
	if win.KeyPressed(KeyEscape) {
		win.uiFocus = 0
	}

	if typed := win.TextInput(); typed != "" {
		// The splice goes through a buffer the window keeps, so typing does
		// not hand the collector a string per keystroke. A field is a handful
		// of characters, and a person types a few a second.
		at := win.uiCursor
		buf := append(win.editBuf[:0], (*text)[:at]...)
		buf = append(buf, typed...)
		buf = append(buf, (*text)[at:]...)
		*text = string(buf)
		win.editBuf = buf
		win.uiCursor = at + len(typed)
		changed = true
		win.uiBlink = 0
	}
	return changed
}
