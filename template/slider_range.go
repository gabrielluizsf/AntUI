package template

// sliderRange is the value a slider shows and the span it sits in, which is
// what the knob's position is worked out from. The value is a pointer because
// the widget writes through it as the knob is dragged, and reading it here is
// what the shared machinery does on the way past.
type sliderRange struct {
	value    *float32
	min, max float32
}

// position is the value across the range as 0..1, clamped at both ends: a
// slider set past its own limits sits at the limit rather than off the track.
func (r sliderRange) position() float32 {
	return min(max((*r.value-r.min)/(r.max-r.min), 0), 1)
}
