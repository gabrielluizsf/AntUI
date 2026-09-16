// Package event is what the templates agree on: the components that exist,
// the things they can do, and the sound bank that hears about them. It has no
// imports of its own, so the visual and the audible halves of a template can
// both speak it without leaning on each other.
//
// The flow is one-way. A component draws itself, and when the user does
// something to it the component reports a single Event: which component it
// was and what was done. The template's sound bank — the second half of the
// template — listens for those events and plays the matching sound. Anything
// can be a sound bank, so a template can be made to say whatever its author
// wants it to say without touching the way it looks.
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
)

// Kind is what happened to a component this frame.
type Kind int

// The things a component can report.
const (
	None   Kind = iota
	Click       // a button was pressed all the way down and released on it
	Toggle      // a checkbox flipped
	Select      // a radio was picked
	Change      // a slider moved
	Type        // text was typed into or deleted from a field
	Focus       // a text field got the keyboard
)

// An Event is one thing a visual component did this frame.
type Event struct {
	Component Component
	Kind      Kind
	Step      int // a value the sound may vary with: where the cursor is, which way a slider went
}

// Nothing is the no-op event a component returns when nothing happened.
var Nothing Event

// Ok reports whether anything happened at all.
func (e Event) Ok() bool { return e.Kind != None }

// Is reports whether the event is the given component doing the given thing.
func (e Event) Is(c Component, k Kind) bool { return e.Component == c && e.Kind == k }

// A SoundBank turns events into sound. It is the hearable half of a template:
// the visual half reports what happened, and the bank decides what it sounds
// like. antui/template/audio ships two of them (Tech and Simple), and any
// type implementing On can play anything at all — or nothing, which is a
// sound bank too.
type SoundBank interface {
	On(e Event)
}
