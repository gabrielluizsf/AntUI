package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// paintCyberKnob draws the thumb: a glow that swells with the hover and press
// state, a neon shell, a dark body and a bright core.
func paintCyberKnob(cv *canvas.Canvas, v vivid, s State, knobX, cy, knobRadius, u int) {
	reach := 22

	if s.Hovered {
		reach = 40
	}

	if s.Active {
		reach = 70
	}

	// Large but soft glow around the knob.
	glow(cv, halo{x: knobX - knobRadius - 3*u, y: cy - knobRadius - 3*u, w: knobRadius*2 + 6*u, h: knobRadius*2 + 6*u, radius: knobRadius + 2*u, colour: v.neon, strength: int(float64(reach) * v.glow), u: u})

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
