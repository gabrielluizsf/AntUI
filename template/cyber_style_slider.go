package template

import "github.com/gabrielluizsf/antui"

// Slider paints a precision control: a thin lit trough, graduations along
// it, and a knob that is a ring around a bright core — a reticle, the same
// language the rest of the theme speaks. The filled run gives off a faint
// glow of its own, so it reads as powered rather than painted.
func (cyberStyle) Slider(win *antui.Window, s State, x, y, w, h int, value float32) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)

	value = min(max(value, 0), 1)

	// Keep the knob completely inside the slider.
	knobRadius := max(h/2, 8*u)

	if w <= knobRadius*2 {
		return
	}

	cy := y + h/2

	// Track runs between the centers of the knob at both extremes.
	trackX := x + knobRadius
	trackW := max(w-2*knobRadius, 1)

	trackH := 8 * u
	trackY := cy - trackH/2

	knobX := trackX + int(value*float32(trackW)+0.5)
	knobX = min(max(knobX, trackX), trackX+trackW)

	paintCyberTrack(cv, v, trackX, trackY, trackW, trackH, knobX, u)
	paintCyberMarks(cv, v, trackX, trackY, trackW, trackH, knobX, u)
	paintCyberTerminals(cv, v, trackX, cy, trackW, u)
	paintCyberKnob(cv, v, s, knobX, cy, knobRadius, u)
}
