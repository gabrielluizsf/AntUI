package template

import "github.com/gabrielluizsf/antui"

// textAreaPress handles a press on the field: a click inside focuses it and
// puts the caret where the pointer pointed, a click outside clears the focus.
// It reports whether the field just gained focus.
func (t *uiTemplate) textAreaPress(win *antui.Window, id uint32, x, y, w, h int, text string, hovered bool) bool {
	if !win.MousePressed(antui.MouseLeft) {
		return false
	}
	switch {
	case hovered:
		justFocused := win.WidgetFocus() != id
		if justFocused {
			win.WidgetFocusSet(id)
		}
		// The caret goes where the pointer pointed, even on a re-click into
		// the field.
		win.WidgetCursorSet(t.style.TextAreaTextPos(win, State{
			Hovered: hovered,
			Focused: win.WidgetFocus() == id,
		}, x, y, w, h, text, win.WidgetCursor(), win.MouseX(), win.MouseY()))
		return justFocused
	case win.WidgetFocus() == id:
		win.WidgetFocusClear(id)
	}
	return false
}
