package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/event"
)

func (t *uiTemplate) Button(win *antui.Window, area canvas.Area, label string) event.Event {
	x, y, w, h := area.X, area.Y, area.Width, area.Height
	id := win.WidgetID("template:button", area, label)
	hovered := win.Hovered(area)
	clicked := win.WidgetClick(id, hovered)
	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)
	keyed := win.WidgetKeyActivate(id)
	if keyed {
		clicked = true
	}
	t.style.Button(ButtonParams{
		Win: win,
		State: State{
			Hovered: hovered,
			Pressed: (win.WidgetActive() == id && hovered) || keyed,
			Focused: focused,
		},
		X: x, Y: y, W: w, H: h,
		Label: label,
	})
	if clicked {
		return t.fire(event.Event{Component: event.Button, Kind: event.Click})
	}
	return event.Nothing
}
