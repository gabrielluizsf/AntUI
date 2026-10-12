package event

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
