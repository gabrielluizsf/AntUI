package antui

import "github.com/gabrielluizsf/antui/canvas"

// sliderDraw paints the track, fill and knob of a slider.
func (win *Window) sliderDraw(s sliderState, value float32) {
	t := (value - s.minValue) / s.span

	trackY := s.y + s.h/2 - 3
	trackX := s.x + s.knobRadius
	knobX := trackX + int(t*float32(s.usable)+0.5)

	win.cv.FillRoundRect(canvas.Round(canvas.Area{X: trackX, Y: trackY, Width: s.usable, Height: 6}, 3), win.theme.Border)
	win.cv.FillRoundRect(canvas.Round(canvas.Area{X: trackX, Y: trackY, Width: knobX - trackX, Height: 6}, 3), win.theme.Accent)

	knob := win.theme.Accent
	if s.hovered || win.uiActive == s.id {
		knob = win.theme.AccentHover
	}
	at := canvas.Circle{X: knobX, Y: s.y + s.h/2, R: s.knobRadius}
	win.cv.FillCircle(at, knob)
	win.cv.FillCircle(canvas.Circle{X: at.X, Y: at.Y, R: at.R - 4}, win.theme.Surface)
	if s.focused {
		win.cv.Circle(canvas.Circle{X: at.X, Y: at.Y, R: at.R + 3}, win.theme.Accent)
	}
}
