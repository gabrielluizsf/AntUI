package template

import ()

func (builtinStyle) Slider(p SliderParams) {
	win := p.Win
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)

	value := min(max(p.Value, 0), 1)

	// Keep the thumb completely inside the control.
	knobRadius := max(p.H/2, 7*u)

	if p.W <= knobRadius*2 {
		return
	}

	cy := p.Y + p.H/2

	// The track goes from the center of the left thumb position
	// to the center of the right thumb position.
	trackX := p.X + knobRadius
	trackW := max(p.W-2*knobRadius, 1)

	trackH := 8 * u
	trackY := cy - trackH/2

	knobX := trackX + int(value*float32(trackW)+0.5)
	knobX = min(max(knobX, trackX), trackX+trackW)

	paintSliderTrack(cv, th, sliderTrack{X: trackX, Y: trackY, Width: trackW, Height: trackH, KnobX: knobX})
	paintSliderThumb(cv, th, sliderThumb{X: knobX, Y: cy, Radius: knobRadius, U: u, State: p.State})
}
