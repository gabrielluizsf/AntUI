package event

// Component names the interactive components every template draws.
type Component int

// The components a template offers. A template that draws all of them answers
// every screen AntUI can ask for; a template may draw a subset and leave the
// rest silent.
const (
	Button Component = iota
	Checkbox
	Radio
	Slider
	TextInput
	Select     // a dropdown picker
	Switch     // a flip toggle, like a checkbox but a thumb
	TextArea   // a multi-line text field
	Progress   // a read-only progress bar
	DatePicker // a calendar the user picks a date on
)
