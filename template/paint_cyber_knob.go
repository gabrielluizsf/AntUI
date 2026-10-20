package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// paintCyberKnob draws the thumb: a solid accent dot separated from the fill
// by a dark ring. Hover and press only grow it a touch — no halo.
func paintCyberKnob(cv *canvas.Canvas, v vivid, s State, knobX, cy, knobRadius, u int) {
	r := knobRadius
	if s.Hovered || s.Active {
		r += 2 * u
	}
	cv.FillCircle(canvas.Circle{X: knobX, Y: cy, R: r}, cyberPalette.bg)
	cv.FillCircle(canvas.Circle{X: knobX, Y: cy, R: r - 2*u}, v.neon)
	if s.Focused {
		cv.Circle(canvas.Circle{X: knobX, Y: cy, R: r + 3*u}, canvas.Fade(v.neon, 120))
	}
}
