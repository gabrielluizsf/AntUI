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
func sliderGeometry(x, y, w, h int, value *float32, minValue, maxValue float32, u int) sliderGeo {
	knobRadius := max(h/2, 6*u)
	span := maxValue - minValue
	usable := max(w-2*knobRadius, 1)
	position := min(max((*value-minValue)/span, 0), 1)
	return sliderGeo{
		span:       span,
		knobRadius: knobRadius,
		usable:     usable,
		trackX:     x + knobRadius,
		knobX:      x + knobRadius + int(position*float32(usable)+0.5),
		knobY:      y + h/2,
	}
}
