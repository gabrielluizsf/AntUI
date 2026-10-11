package template

import "github.com/gabrielluizsf/antui"

// editTextArea is the multi-line cousin of WidgetEdit: it applies this frame's
// keys and typing to the focused field. maxRef is the wrap width in reference
// pixels (the width the text is measured at scale 1).
func (t *uiTemplate) editTextArea(win *antui.Window, text *string, maxRef int) bool {
	win.WidgetCursorSet(clampInt(win.WidgetCursor(), 0, len(*text)))
	changed := false

	if win.KeyPressed(antui.KeyBackspace) && win.WidgetCursor() > 0 {
		start := utf8PrevCursor(*text, win.WidgetCursor())
		*text = (*text)[:start] + (*text)[win.WidgetCursor():]
		win.WidgetCursorSet(start)
		changed = true
	}
	if win.KeyPressed(antui.KeyDelete) && win.WidgetCursor() < len(*text) {
		end := utf8NextCursor(*text, win.WidgetCursor())
		*text = (*text)[:win.WidgetCursor()] + (*text)[end:]
		changed = true
	}
	if win.KeyPressed(antui.KeyLeft) && win.WidgetCursor() > 0 {
		win.WidgetCursorSet(utf8PrevCursor(*text, win.WidgetCursor()))
	}
	if win.KeyPressed(antui.KeyRight) && win.WidgetCursor() < len(*text) {
		win.WidgetCursorSet(utf8NextCursor(*text, win.WidgetCursor()))
	}
	if win.KeyPressed(antui.KeyEnter) {
		c := win.WidgetCursor()
		*text = (*text)[:c] + "\n" + (*text)[c:]
		win.WidgetCursorSet(utf8NextCursor(*text, c))
		changed = true
	}

	// Navigation that depends on the wrapped layout: Up/Down move between
	// visual lines, Home/End snap to the edges of the current visual line.
	lines := textLines(*text, maxRef)
	li, in := cursorLine(lines, win.WidgetCursor())
	if win.KeyPressed(antui.KeyUp) {
		win.WidgetCursorSet(moveLineUp(*text, lines, li, in))
	}
	if win.KeyPressed(antui.KeyDown) {
		win.WidgetCursorSet(moveLineDown(*text, lines, li, in))
	}
	if win.KeyPressed(antui.KeyHome) {
		win.WidgetCursorSet(lines[li][0])
	}
	if win.KeyPressed(antui.KeyEnd) {
		win.WidgetCursorSet(lines[li][1])
	}

	if typed := win.TextInput(); typed != "" {
		c := win.WidgetCursor()
		*text = (*text)[:c] + typed + (*text)[c:]
		win.WidgetCursorSet(c + len(typed))
		changed = true
	}
	return changed
}
