package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/event"
)

func (t *uiTemplate) Slider(
	win *antui.Window,
	x, y, w, h int,
	value *float32,
	minValue, maxValue float32,
) event.Event {
	if value == nil || w <= 0 || h <= 0 || maxValue <= minValue {
		return event.Nothing
	}

	id := win.WidgetID("template:slider", x, y, w, h, "")

	u := Scale(win)

	// The knob is the only clickable part of the slider.
	knobRadius := max(h/2, 6*u)

	span := maxValue - minValue

	// The usable movement is from the center of the left knob
	// to the center of the right knob.
	usable := max(w-2*knobRadius, 1)

	// Convert the current value to a normalized position.
	position := (*value - minValue) / span
	position = min(max(position, 0), 1)

	// Current knob center.
	knobX := x + knobRadius + int(position*float32(usable)+0.5)
	knobY := y + h/2

	mouseX := win.MouseX()
	mouseY := win.MouseY()

	// ---------------------------------------------------------------------
	// Hit testing
	// ---------------------------------------------------------------------

	dx := mouseX - knobX
	dy := mouseY - knobY

	knobHit := dx*dx+dy*dy <= knobRadius*knobRadius

	if knobHit {
		win.WidgetHot(id)
	}

	// ---------------------------------------------------------------------
	// Start dragging
	// ---------------------------------------------------------------------

	if knobHit && win.MousePressed(antui.MouseLeft) {
		win.WidgetStartPress(id)
		win.WidgetFocusSet(id)

		// Preserve the exact point where the mouse grabbed the knob.
		//
		// Example:
		//     knob center = 200
		//     mouse       = 197
		//
		// We remember -3 so the knob does not jump to the mouse center.
		t.sliderDragOffset = mouseX - knobX
	}

	active := win.WidgetActive() == id

	// ---------------------------------------------------------------------
	// End dragging
	// ---------------------------------------------------------------------

	if active && !win.MouseDown(antui.MouseLeft) {
		win.WidgetEndPress(id)
		active = false
		t.sliderDragOffset = 0
	}

	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)

	changed := false

	// ---------------------------------------------------------------------
	// Keyboard
	// ---------------------------------------------------------------------

	if focused {
		step := span / 10

		if win.KeyPressed(antui.KeyLeft) {
			wanted := *value - step
			wanted = min(max(wanted, minValue), maxValue)

			if wanted != *value {
				*value = wanted
				changed = true
			}
		}

		if win.KeyPressed(antui.KeyRight) {
			wanted := *value + step
			wanted = min(max(wanted, minValue), maxValue)

			if wanted != *value {
				*value = wanted
				changed = true
			}
		}
	}

	// ---------------------------------------------------------------------
	// Mouse dragging
	// ---------------------------------------------------------------------

	if active {
		// Compensate for the point where the mouse grabbed the knob.
		//
		// Without this:
		//
		//     mouse ────────●
		//                    ↑
		//               knob jumps
		//
		// With the offset:
		//
		//     mouse ────●
		//               ↑
		//          knob follows naturally
		//
		adjustedMouseX := mouseX - t.sliderDragOffset

		position := float32(
			adjustedMouseX-(x+knobRadius),
		) / float32(usable)

		position = min(max(position, 0), 1)

		wanted := minValue + position*span
		wanted = min(max(wanted, minValue), maxValue)

		if wanted != *value {
			*value = wanted
			changed = true
		}
	}

	// ---------------------------------------------------------------------
	// Final clamp
	// ---------------------------------------------------------------------

	*value = min(max(*value, minValue), maxValue)

	// ---------------------------------------------------------------------
	// Draw
	// ---------------------------------------------------------------------

	visualValue := (*value - minValue) / span
	visualValue = min(max(visualValue, 0), 1)

	t.style.Slider(
		win,
		State{
			Hovered: knobHit,
			Active:  active,
			Focused: focused,
		},
		x, y, w, h,
		visualValue,
	)

	if changed {
		return t.fire(event.Event{
			Component: event.Slider,
			Kind:      event.Change,
		})
	}

	return event.Nothing
}
