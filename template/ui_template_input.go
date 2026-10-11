package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/event"
)

func (t *uiTemplate) Input(win *antui.Window, x, y, w, h int, text *string) event.Event {
	if text == nil {
		return event.Nothing
	}
	id := win.WidgetID("template:input", x, y, w, h, "")
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
			win.WidgetCursorSet(t.style.InputTextPos(win, State{
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
		changed = win.WidgetEdit(text)
	}

	t.style.Input(win, State{
		Hovered: hovered,
		Focused: focused,
	}, x, y, w, h, *text, win.WidgetCursor(), focused && win.WidgetCaret())

	switch {
	case justFocused:
		return t.fire(event.Event{Component: event.TextInput, Kind: event.Focus})
	case changed:
		return t.fire(event.Event{
			Component: event.TextInput,
			Kind:      event.Type,
			Step:      win.WidgetCursor(),
		})
	}
	return event.Nothing
}
