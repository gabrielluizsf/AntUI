package template

// sliderGeo is the slider's clickable geometry, worked out once from the box
// the field gave it.
type sliderGeo struct {
	span       float32 // the range a position moves the value across
	knobRadius int     // half the knob, the clickable radius
	usable     int     // the distance the knob can travel
	trackX     int     // where the usable run starts, x + knobRadius
	knobX      int     // the knob's center
	knobY      int
}

// sliderGeometry derives the knob's geometry from the box, the value range and
// the scale, with the current value normalized into a knob position.
func sliderGeometry(x, y, w, h int, r sliderRange, u int) sliderGeo {
	knobRadius := max(h/2, 6*u)
	usable := max(w-2*knobRadius, 1)
	position := r.position()
	return sliderGeo{
		span:       r.max - r.min,
		knobRadius: knobRadius,
		usable:     usable,
		trackX:     x + knobRadius,
		knobX:      x + knobRadius + int(position*float32(usable)+0.5),
		knobY:      y + h/2,
	}
}
