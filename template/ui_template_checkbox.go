package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/event"
)

func (t *uiTemplate) Checkbox(win *antui.Window, x, y int, label string, value *bool) event.Event {
	if value == nil {
		return event.Nothing
	}
	u := Scale(win)
	box := 18 * u
	w := box + 8*u + textWidth(u, label)
	id := win.WidgetID("template:checkbox", x, y, box, box, label)
	hovered := win.Hovered(x, y, w, box)
	clicked := win.WidgetClick(id, hovered)
	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)
	if win.WidgetKeyActivate(id) {
		clicked = true
	}
	changed := false
	if clicked {
		*value = !*value
		changed = true
	}
	t.style.Checkbox(win, State{Hovered: hovered, Focused: focused, On: *value}, x, y, label, *value)
	if changed {
		return t.fire(event.Event{Component: event.Checkbox, Kind: event.Toggle})
	}
	return event.Nothing
}
