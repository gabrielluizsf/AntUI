package antui

// Gestures is what the fingers did this frame. Every field is about this
// frame alone: a Tap is true on the frame the finger lifted and false on the
// next one, the same way MousePressed is.
//
// It is touch only. A mouse has no gestures — it has buttons and a wheel,
// which are reported as themselves — so on a desktop this is always empty.
type Gestures struct {
	// X and Y are where the gesture is: the finger, or the point between two
	// of them for a pinch.
	X, Y int

	// Tap is a finger that landed and lifted without wandering. Count is 1
	// for a single tap, 2 for a double, and keeps going for as long as the
	// taps keep coming quickly enough — so a triple tap is Count 3, not
	// three separate taps.
	Tap   bool
	Count int

	// LongPress fires once, while the finger is still down, as soon as it
	// has been still for long enough. It does not wait for the lift, because
	// by then the user has been waiting for half a second with no answer.
	LongPress bool

	// Drag is a finger that has moved past the slop. DX and DY are this
	// frame's movement.
	Drag   bool
	DX, DY int

	// Fling is a drag that was let go while still moving. VX and VY are how
	// fast, in pixels a second.
	Fling  bool
	VX, VY float64

	// Pinch and Rotate need two fingers. Scale is the change this frame —
	// above 1 for fingers moving apart — and Angle is the change in radians,
	// positive clockwise.
	Pinch  bool
	Scale  float64
	Rotate bool
	Angle  float64
}
