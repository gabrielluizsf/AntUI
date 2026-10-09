package css

// styleKey is what decides an answer: the widget, the state it is in and the
// viewport the cascade read. Media queries test the viewport, so two windows
// of the same width but different heights are two different keys.
type styleKey struct {
	tag    string
	state  State
	width  int
	height int
}
