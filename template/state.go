package template

// State is everything the interaction layer knows about a widget this frame,
// handed to a style so it can paint the right state with no guessing.
type State struct {
	Hovered bool // the pointer is over the widget
	Pressed bool // the widget is being pressed right now
	Active  bool // the widget is the one being dragged (a slider)
	Focused bool // the widget has the keyboard; Tab walks between these
	On      bool // for checkboxes and radios: the value it shows
}
