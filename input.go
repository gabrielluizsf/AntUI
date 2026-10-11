package antui

// Input draws a single-line text field and reports whether the text changed
// this frame.
func (win *Window) Input(x, y, w, h int, text *string) bool {
	if text == nil {
		return false
	}
	s := inputState{x: x, y: y, w: w, h: h, id: widgetID("input", x, y, w, h, "")}
	s.hovered = win.Hovered(x, y, w, h)
	if s.hovered {
		win.uiHot = s.id
	}

	if win.MousePressed(MouseLeft) {
		switch {
		case s.hovered:
			if win.uiFocus != s.id {
				win.uiFocus = s.id
				win.uiCursor = len(*text)
			}
			win.uiBlink = 0
		case win.uiFocus == s.id:
			win.uiFocus = 0
		}
	}
	s.focused = win.uiFocus == s.id
	win.uiTab = append(win.uiTab, s.id)

	changed := false
	if s.focused {
		changed = win.editText(text)
	}

	win.drawInput(s, text)
	return changed
}
