package canvas

// The shapes a drawing is made of, as values rather than as loose numbers, so
// a rectangle is passed around as a rectangle and not as four ints the reader
// has to guess the order of. Each is a plain value: copying one costs nothing
// and nothing here holds a reference to a canvas.

// Rounded is an [Area] with its corners rounded — by RX across and RY down, so
// a square corner is RX and RY both zero and a circle is the two equal at half
// the shorter side. A RY of zero means the same as RX, the way the CSS radii do.
type Rounded struct {
	Area
	RX, RY int
}

// Round is [Rounded] with the two radii equal, which is what almost every
// rounded rectangle in a UI actually is.
func Round(a Area, r int) Rounded { return Rounded{Area: a, RX: r, RY: r} }

// vertical reports the radius down the box, which is RX when the caller left
// RY empty.
func (r Rounded) vertical() int {
	if r.RY == 0 {
		return r.RX
	}
	return r.RY
}

// Circle is a circle: where its middle is and how far it reaches.
type Circle struct {
	X, Y, R int
}

// Ellipse is an ellipse with a radius each way, so a circle is an ellipse whose
// two are equal.
type Ellipse struct {
	X, Y   int
	RX, RY int
}
