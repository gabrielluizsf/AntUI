package css

// BackPosition is one axis of a background-position, or a gradient's centre:
// an edge to anchor the picture to and an offset away from it. The offset is
// a length in reference pixels or a percentage of the side left over — the
// positioning area minus the image, as CSS measures it.
type BackPosition struct {
	Edge uint8 // one of the BackPos* constants
	Off  Length
}

const (
	BackPosStart  uint8 = iota // the left/top edge
	BackPosMiddle              // the centre
	BackPosEnd                 // the right/bottom edge
	BackPosOffset              // a bare length/percentage, no keyword
)
