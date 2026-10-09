package css

// layoutStyle carries the layout fields of a Style: what kind of box it is,
// where it sits, and how it behaves around its neighbours.
type layoutStyle struct {
	Display       uint8 // one of the Display* constants
	Position      uint8 // one of the Position* constants
	Top           Length
	Right         Length
	Bottom        Length
	Left          Length
	ZIndex        int
	Overflow      [2]uint8 // x, y; one of the Overflow* constants
	Visibility    uint8    // one of the Visibility* constants
	PointerEvents uint8    // one of the PointerEvents* constants
	Cursor        uint8    // one of the Cursor* constants
}

// OutOfFlow reports whether the box is positioned outside the normal flow.
func (s Style) OutOfFlow() bool {
	return s.Position == PositionAbsolute || s.Position == PositionFixed
}
