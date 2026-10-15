package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/event"
)

// Switch draws a flip toggle: a track with a thumb that slides between off
// and on. It flips value when clicked and reports whether it flipped.
func (t *uiTemplate) Switch(win *antui.Window, area canvas.Area, params antui.CheckboxParams) event.Event {
	label := params.Label
	value := params.Value
	x, y := area.X, area.Y
	if value == nil {
		return event.Nothing
	}
	u := Scale(win)
	trackW := switchTrackW * u
	trackH := switchTrackH * u
	w := trackW + labelGap*u + textWidth(u, label)
	id := win.WidgetID("template:switch", canvas.Area{X: x, Y: y, Width: trackW, Height: trackH}, label)
	hovered := win.Hovered(canvas.Area{X: x, Y: y, Width: w, Height: trackH})
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
	t.style.Switch(SwitchParams{
		Win:   win,
		State: State{Hovered: hovered, Focused: focused, On: *value},
		X:     x, Y: y,
		Label: label,
		On:    *value,
	})
	if changed {
		return t.fire(event.Event{Component: event.Switch, Kind: event.Toggle})
	}
	return event.Nothing
}
