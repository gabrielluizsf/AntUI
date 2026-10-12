package template

import (
	"github.com/gabrielluizsf/antui"
)

// sliderKeyboard moves the value by a tenth of the span per arrow key press
// and reports whether it changed.
func (t *uiTemplate) sliderKeyboard(win *antui.Window, g sliderGeo, value *float32, minValue, maxValue float32) bool {
	changed := false
	step := g.span / 10
	if win.KeyPressed(antui.KeyLeft) {
		wanted := min(max(*value-step, minValue), maxValue)
		if wanted != *value {
			*value = wanted
			changed = true
		}
	}
	if win.KeyPressed(antui.KeyRight) {
		wanted := min(max(*value+step, minValue), maxValue)
		if wanted != *value {
			*value = wanted
			changed = true
		}
	}
	return changed
}
