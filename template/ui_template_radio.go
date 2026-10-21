package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/event"
)

func (t *uiTemplate) Radio(win *antui.Window, area canvas.Area, params antui.RadioParams) event.Event {
	label := params.Label
	value := params.Value
	option := params.Option
	x, y := area.X, area.Y
	if value == nil {
		return event.Nothing
	}
	u := Scale(win)
	size := boxSize * u
	w := size + labelGap*u + textWidth(u, label)
	id := win.WidgetID("template:radio", canvas.Area{X: x, Y: y, Width: size, Height: option}, label)
	hovered := win.Hovered(canvas.Area{X: x, Y: y, Width: w, Height: size})
	clicked := win.WidgetClick(id, hovered)
	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)
	if win.WidgetKeyActivate(id) {
		clicked = true
	}
	selected := *value == option

	if clicked {
		if selected {
			*value = antui.NoOption
			selected = false
		} else {
			*value = option
			selected = true
		}
	}
	t.style.Radio(RadioParams{
		Win:   win,
		State: State{Hovered: hovered, Focused: focused, On: selected},
		X:     x, Y: y,
		Label: label,
		On:    selected,
	})
	if clicked {
		return t.fire(event.Event{Component: event.Radio, Kind: event.Pick})
	}
	return event.Nothing
}
