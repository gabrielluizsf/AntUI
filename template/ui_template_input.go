package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/event"
)

func (t *uiTemplate) Input(win *antui.Window, x, y, w, h int, text *string) event.Event {
	if text == nil {
		return event.Nothing
	}
	id := win.WidgetID("template:input", canvas.Area{X: x, Y: y, Width: w, Height: h}, "")
	hovered := win.Hovered(canvas.Area{X: x, Y: y, Width: w, Height: h})
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
			win.WidgetCursorSet(t.style.InputTextPos(InputTextPosParams{
				Win: win,
				State: State{
					Hovered: hovered,
					Focused: win.WidgetFocus() == id,
				},
				X:      x,
				Y:      y,
				W:      w,
				H:      h,
				Text:   *text,
				Cursor: win.WidgetCursor(),
				MouseX: win.MouseX(),
				MouseY: win.MouseY(),
			}))
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

	t.style.Input(InputParams{
		Win: win,
		State: State{
			Hovered: hovered,
			Focused: focused,
		},
		X:      x,
		Y:      y,
		W:      w,
		H:      h,
		Text:   *text,
		Cursor: win.WidgetCursor(),
		Caret:  focused && win.WidgetCaret(),
	})

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
