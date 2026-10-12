package template

import (
	"github.com/gabrielluizsf/antui"
)

// sliderDrag follows the mouse while a drag is active, compensating for the
// point where the mouse grabbed the knob, and reports whether the value
// moved.
func (t *uiTemplate) sliderDrag(win *antui.Window, g sliderGeo, mouseX int, value *float32, minValue, maxValue float32) bool {
	adjustedMouseX := mouseX - t.sliderDragOffset
	position := float32(adjustedMouseX-g.trackX) / float32(g.usable)
	position = min(max(position, 0), 1)
	wanted := minValue + position*g.span
	wanted = min(max(wanted, minValue), maxValue)
	if wanted != *value {
		*value = wanted
		return true
	}
	return false
}
