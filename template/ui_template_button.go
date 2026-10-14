package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/event"
)

func (t *uiTemplate) Button(win *antui.Window, x, y, w, h int, label string) event.Event {
	id := win.WidgetID("template:button", canvas.Area{X: x, Y: y, Width: w, Height: h}, label)
	hovered := win.Hovered(canvas.Area{X: x, Y: y, Width: w, Height: h})
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
