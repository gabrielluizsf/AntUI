package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// paintSliderThumb draws the knob, its inner surface and the focus ring.
func paintSliderThumb(cv *canvas.Canvas, th *antui.Theme, s State, knobX, cy, knobRadius, u int) {
	knob := th.Accent
	if s.Active || s.Hovered {
		knob = th.AccentHover
	}

	// Outer thumb.
	cv.FillCircle(knobX, cy, knobRadius, knob)

	// Inner surface creates the visual "thumb" ring.
	innerRadius := max(knobRadius-4*u, 2*u)
	cv.FillCircle(knobX, cy, innerRadius, th.Surface)

	// Focus ring.
	if s.Focused {
		cv.Circle(knobX, cy, knobRadius+3*u, th.Accent)
	}
}
