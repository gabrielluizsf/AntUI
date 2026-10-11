package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/event"
)

// TextArea draws a multi-line text field: Enter inserts a line break, the
// arrows and Home/End move the cursor between and along wrapped lines, and
// text flows onto as many lines as will fit in the box. It reports when the
// text changed or the field gained focus.
func (t *uiTemplate) TextArea(win *antui.Window, x, y, w, h int, text *string) event.Event {
	if text == nil {
		return event.Nothing
	}
	id := win.WidgetID("template:textarea", x, y, w, h, "")
	hovered := win.Hovered(x, y, w, h)
	if hovered {
		win.WidgetHot(id)
	}

	justFocused := false
	if win.MousePressed(antui.MouseLeft) {
		switch {
		case hovered:
			if win.WidgetFocus() != id {
				win.WidgetFocusSet(id)
				justFocused = true
			}
			// The caret goes where the pointer pointed, even on a re-click
			// into the field.
			win.WidgetCursorSet(t.style.TextAreaTextPos(win, State{
				Hovered: hovered,
				Focused: win.WidgetFocus() == id,
			}, x, y, w, h, *text, win.WidgetCursor(), win.MouseX(), win.MouseY()))
		case win.WidgetFocus() == id:
			win.WidgetFocusClear(id)
		}
	}
	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)

	changed := false
	if focused {
		u := Scale(win)
		maxRef := max(w/(8*u), 12)
		changed = t.editTextArea(win, text, maxRef)
	}

	t.style.TextArea(win, State{
		Hovered: hovered,
		Focused: focused,
	}, x, y, w, h, *text, win.WidgetCursor(), focused && win.WidgetCaret())

	switch {
	case justFocused:
		return t.fire(event.Event{Component: event.TextArea, Kind: event.Focus})
	case changed:
		return t.fire(event.Event{
			Component: event.TextArea,
			Kind:      event.Type,
			Step:      win.WidgetCursor(),
		})
	}
	return event.Nothing
}
