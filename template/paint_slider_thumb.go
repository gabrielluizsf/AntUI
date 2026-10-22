package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// sliderThumb is where the knob sits, how big it is and how it is being
// touched right now, at the scale the window draws this frame.
type sliderThumb struct {
	X, Y, Radius int
	U            int
	State        State
}

// paintSliderThumb draws the knob, its drop shadow, its inner surface and the
// focus ring.
func paintSliderThumb(cv *canvas.Canvas, th *antui.Theme, k sliderThumb) {
	s := k.State
	knob := th.Accent
	if s.Active || s.Hovered {
		knob = th.AccentHover
	}

	// A thumb is the one part of the control a finger actually lands on, so
	// it is the one part that is lifted off the track rather than flush
	// with it.
	cv.BoxShadow(canvas.BoxShadowParams{
		X: k.X - k.Radius, Y: k.Y - k.Radius, Width: k.Radius * 2, Height: k.Radius * 2,
		RX: k.Radius, RY: k.Radius, DY: k.U, Blur: 3 * k.U, Color: th.Shadow,
	})
	cv.FillCircle(canvas.Circle{X: k.X, Y: k.Y, R: k.Radius}, knob)

	// Inner surface creates the visual "thumb" ring.
	innerRadius := max(k.Radius-4*k.U, 2*k.U)
	cv.FillCircle(canvas.Circle{X: k.X, Y: k.Y, R: innerRadius}, th.Surface)

	if s.Focused {
		cv.Circle(canvas.Circle{X: k.X, Y: k.Y, R: k.Radius + 3*k.U}, th.Accent)
	}
}
