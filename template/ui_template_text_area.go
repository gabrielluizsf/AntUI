package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/event"
)

// TextArea draws a multi-line text field: Enter inserts a line break, the
// arrows and Home/End move the cursor between and along wrapped lines, and
// text flows onto as many lines as will fit in the box. It reports when the
// text changed or the field gained focus.
func (t *uiTemplate) TextArea(win *antui.Window, area canvas.Area, text *string) event.Event {
	x, y, w, h := area.X, area.Y, area.Width, area.Height
	if text == nil {
		return event.Nothing
	}
	id := win.WidgetID("template:textarea", canvas.Area{X: x, Y: y, Width: w, Height: h}, "")
	hovered := win.Hovered(canvas.Area{X: x, Y: y, Width: w, Height: h})
	if hovered {
		win.WidgetHot(id)
	}

	justFocused := t.textAreaPress(win, id, x, y, w, h, *text, hovered)
	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)

	changed := false
	if focused {
		u := Scale(win)
		changed = t.editTextArea(win, text, max(w/(8*u), 12))
	}

	t.style.TextArea(TextAreaParams{
		Win:   win,
		State: State{Hovered: hovered, Focused: focused},
		X:     x, Y: y, W: w, H: h,
		Text:   *text,
		Cursor: win.WidgetCursor(),
		Caret:  focused && win.WidgetCaret(),
	})

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
