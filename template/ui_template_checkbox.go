package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/event"
)

func (t *uiTemplate) Checkbox(win *antui.Window, x, y int, label string, value *bool) event.Event {
	if value == nil {
		return event.Nothing
	}
	u := Scale(win)
	box := boxSize * u
	w := box + labelGap*u + textWidth(u, label)
	id := win.WidgetID("template:checkbox", canvas.Area{X: x, Y: y, Width: box, Height: box}, label)
	hovered := win.Hovered(canvas.Area{X: x, Y: y, Width: w, Height: box})
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
	t.style.Checkbox(CheckboxParams{
		Win:   win,
		State: State{Hovered: hovered, Focused: focused, On: *value},
		X:     x, Y: y,
		Label: label,
		On:    *value,
	})
	if changed {
		return t.fire(event.Event{Component: event.Checkbox, Kind: event.Toggle})
	}
	return event.Nothing
}
