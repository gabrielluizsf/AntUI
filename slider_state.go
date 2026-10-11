package antui

// sliderState is the resolved layout of a slider for one frame.
type sliderState struct {
	x, y, w, h int
	id         uint32
	hovered    bool
	focused    bool
	knobRadius int
	usable     int
	span       float32
	minValue   float32
	maxValue   float32
}
