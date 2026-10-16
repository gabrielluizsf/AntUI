package svg

// The size a picture drawn into a box that was not given one comes out at,
// which is the default a replaced element takes in CSS: three hundred across
// and a hundred and fifty down.
const (
	defaultPictureW = 300
	defaultPictureH = 150
)

// autoBox is the box with the sides the drawing left open taken from the
// picture that arrived: its own size for a box left open across the board, and
// its proportion for a box one side of which was written — a wide photograph
// asked to be four across comes out two down, as it was encoded. A picture
// that answers with no size of its own, which no format this reads does today,
// takes the default above rather than a side of nothing.
func autoBox(box [4]float64, autoW, autoH bool, ownW, ownH float64) [4]float64 {
	if !autoW && !autoH {
		return box
	}
	if ownW <= 0 || ownH <= 0 {
		ownW, ownH = defaultPictureW, defaultPictureH
	}
	switch {
	case autoW && autoH:
		box[2], box[3] = ownW, ownH
	case autoW:
		box[2] = box[3] * ownW / ownH
	default:
		box[3] = box[2] * ownH / ownW
	}
	return box
}
