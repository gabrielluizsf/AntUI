package antui

// Touch is one finger on the screen, from the moment it lands to the moment
// it lifts.
//
// It is not a mouse and is not reported as one. A mouse has one position and
// a set of buttons; fingers arrive several at a time, each with its own
// history, and a program that wants to pinch or to let two people play on
// one screen needs to tell them apart. For everything simpler, the first
// finger down is also reported as the left mouse button, so a program
// written for a desktop works on a phone without knowing.
type Touch struct {
	// ID stays with this finger for as long as it is down, and may be reused
	// by a later one. It is not an index: a finger's position in
	// [Window.Touches] changes as others come and go.
	ID int

	// X and Y are where it is now, in canvas pixels.
	X, Y int
	// DX and DY are how far it moved this frame.
	DX, DY int
	// StartX and StartY are where it landed, which is what a swipe is
	// measured from.
	StartX, StartY int

	// Pressure is 0 to 1 on hardware that measures it. On hardware that does
	// not it is 1 the whole time the finger is down, so it cannot be used to
	// find out whether the device reports pressure at all.
	Pressure float32
	// Size is how much of the screen the touch covers, 0 to 1. It is how a
	// thumb is told from a fingertip.
	Size float32
	// Tool is what is doing the touching.
	Tool Tool

	// Began is set for the one frame the finger landed on, and Ended for the
	// one frame it lifted on. A touch is still in [Window.Touches] on the
	// frame it ended, so nothing has to be caught between frames.
	Began bool
	Ended bool
	// Cancelled says the system took the touch away rather than the user
	// lifting it — a notification pulled down over the app, a gesture the
	// system claimed. A tap must not be counted from a cancelled touch.
	Cancelled bool
}
