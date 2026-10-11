package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/event"
)

func (t *uiTemplate) Button(win *antui.Window, x, y, w, h int, label string) event.Event {
	id := win.WidgetID("template:button", x, y, w, h, label)
	hovered := win.Hovered(x, y, w, h)
	clicked := win.WidgetClick(id, hovered)
	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)
	keyed := win.WidgetKeyActivate(id)
	if keyed {
		clicked = true
	}
	t.style.Button(win, State{
		Hovered: hovered,
		Pressed: (win.WidgetActive() == id && hovered) || keyed,
		Focused: focused,
	}, x, y, w, h, label)
	if clicked {
		return t.fire(event.Event{Component: event.Button, Kind: event.Click})
	}
	return event.Nothing
}
