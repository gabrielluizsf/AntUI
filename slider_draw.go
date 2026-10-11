package antui

// sliderDraw paints the track, fill and knob of a slider.
func (win *Window) sliderDraw(s sliderState, value float32) {
	t := (value - s.minValue) / s.span

	trackY := s.y + s.h/2 - 3
	knobX := s.x + s.knobRadius + int(t*float32(s.usable)+0.5)

	win.cv.FillRoundRect(s.x+s.knobRadius, trackY, s.usable, 6, 3, win.theme.Border)
	win.cv.FillRoundRect(s.x+s.knobRadius, trackY, knobX-(s.x+s.knobRadius), 6, 3, win.theme.Accent)

	knob := win.theme.Accent
	if s.hovered || win.uiActive == s.id {
		knob = win.theme.AccentHover
	}
	win.cv.FillCircle(knobX, s.y+s.h/2, s.knobRadius, knob)
	win.cv.FillCircle(knobX, s.y+s.h/2, s.knobRadius-4, win.theme.Surface)
	if s.focused {
		win.cv.Circle(knobX, s.y+s.h/2, s.knobRadius+3, win.theme.Accent)
	}
}
