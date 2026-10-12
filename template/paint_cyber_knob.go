package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// paintCyberKnob draws the thumb: a glow that swells with the hover and press
// state, a neon shell, a dark body and a bright core.
func paintCyberKnob(cv *canvas.Canvas, v vivid, s State, knobX, cy, knobRadius, u int) {
	halo := 22

	if s.Hovered {
		halo = 40
	}

	if s.Active {
		halo = 70
	}

	// Large but soft glow around the knob.
	glow(cv, knobX-knobRadius-3*u, cy-knobRadius-3*u, knobRadius*2+6*u, knobRadius*2+6*u, knobRadius+2*u, v.neon, int(float64(halo)*v.glow), u)

	// Outer neon shell.
	cv.FillCircle(knobX, cy, knobRadius+2*u, v.neon)

	// Dark body.
	cv.FillCircle(knobX, cy, knobRadius, cyberPalette.panelAct)

	// Bright inner ring.
	cv.Circle(knobX, cy, knobRadius-2*u, v.neon)

	// Inner core.
	cv.FillCircle(knobX, cy, 3*u, v.neon)

	// Focus gives the thumb an additional outer ring.
	if s.Focused {
		cv.Circle(knobX, cy, knobRadius+5*u, v.text)
	}
}
