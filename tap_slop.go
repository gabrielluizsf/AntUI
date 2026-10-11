package antui

// The thresholds a gesture is recognised against, in points rather than
// pixels so that they mean the same thing on a dense screen as on a coarse
// one. They are the platform's own numbers: a gesture that behaves
// differently from every other app on the device is a gesture the user gets
// wrong, however well it is implemented.
const (
	// TapSlop is how far a finger may wander and still be a tap rather than
	// a drag.
	TapSlop = 8.0
	// LongPressTime is how long a finger must stay put to be a long press.
	LongPressTime = 0.5
	// DoubleTapTime and DoubleTapSlop bound the gap between two taps for
	// them to count as one double tap.
	DoubleTapTime = 0.3
	DoubleTapSlop = 32.0
	// FlingSpeed is the slowest lift that still counts as a throw, in points
	// a second.
	FlingSpeed = 50.0
)
