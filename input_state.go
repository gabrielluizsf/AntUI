package antui

// inputState is the resolved layout of a text field for one frame.
type inputState struct {
	x, y, w, h int
	id         uint32
	hovered    bool
	focused    bool
}
