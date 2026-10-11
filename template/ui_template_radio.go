package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/event"
)

func (t *uiTemplate) Radio(win *antui.Window, x, y int, label string, value *int, option int) event.Event {
	if value == nil {
		return event.Nothing
	}
	u := Scale(win)
	size := 18 * u
	w := size + 8*u + textWidth(u, label)
	id := win.WidgetID("template:radio", x, y, size, option, label)
	hovered := win.Hovered(x, y, w, size)
	clicked := win.WidgetClick(id, hovered)
	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)
	if win.WidgetKeyActivate(id) {
		clicked = true
	}
	selected := *value == option

	if clicked && !selected {
		*value = option
		selected = true
	} else {
		clicked = false
	}
	t.style.Radio(win, State{Hovered: hovered, Focused: focused, On: selected}, x, y, label, selected)
	if clicked {
		return t.fire(event.Event{Component: event.Radio, Kind: event.Pick})
	}
	return event.Nothing
}
