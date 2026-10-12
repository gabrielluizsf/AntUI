package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/event"
)

func (t *uiTemplate) Slider(win *antui.Window, x, y, w, h int, value *float32, minValue, maxValue float32) event.Event {
	if value == nil || w <= 0 || h <= 0 || maxValue <= minValue {
		return event.Nothing
	}

	id := win.WidgetID("template:slider", x, y, w, h, "")
	g := sliderGeometry(x, y, w, h, value, minValue, maxValue, Scale(win))
	mouseX, mouseY := win.MouseX(), win.MouseY()

	// The knob is the only clickable part; a hit warms it.
	dx, dy := mouseX-g.knobX, mouseY-g.knobY
	knobHit := dx*dx+dy*dy <= g.knobRadius*g.knobRadius
	if knobHit {
		win.WidgetHot(id)
	}

	// The press is remembered where the mouse grabbed the knob, so the knob
	// follows the pointer instead of jumping to its center.
	if knobHit && win.MousePressed(antui.MouseLeft) {
		win.WidgetStartPress(id)
		win.WidgetFocusSet(id)
		t.sliderDragOffset = mouseX - g.knobX
	}

	active := win.WidgetActive() == id
	if active && !win.MouseDown(antui.MouseLeft) {
		win.WidgetEndPress(id)
		active = false
		t.sliderDragOffset = 0
	}

	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)

	changed := false
	if focused {
		changed = t.sliderKeyboard(win, g, value, minValue, maxValue)
	}
	if active {
		changed = t.sliderDrag(win, g, mouseX, value, minValue, maxValue) || changed
	}

	*value = min(max(*value, minValue), maxValue)

	visualValue := min(max((*value-minValue)/g.span, 0), 1)
	t.style.Slider(win, State{Hovered: knobHit, Active: active, Focused: focused}, x, y, w, h, visualValue)

	if changed {
		return t.fire(event.Event{Component: event.Slider, Kind: event.Change})
	}
	return event.Nothing
}
