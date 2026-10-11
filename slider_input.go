package antui

// sliderInput applies keyboard stepping and pointer dragging, reporting
// whether the value changed.
func (win *Window) sliderInput(s sliderState, value *float32) bool {
	changed := false
	if s.focused {
		step := s.span / 10
		if win.KeyPressed(KeyLeft) {
			wanted := *value - step
			if wanted != *value {
				*value = min(max(wanted, s.minValue), s.maxValue)
				changed = true
				win.uiBlink = 0
			}
		}
		if win.KeyPressed(KeyRight) {
			wanted := *value + step
			if wanted != *value {
				*value = min(max(wanted, s.minValue), s.maxValue)
				changed = true
				win.uiBlink = 0
			}
		}
	}
	if win.uiActive == s.id {
		position := float32(win.mouseX-(s.x+s.knobRadius)) / float32(s.usable)
		position = min(max(position, 0), 1)
		if wanted := s.minValue + position*s.span; wanted != *value {
			*value = wanted
			changed = true
		}
	}
	return changed
}
