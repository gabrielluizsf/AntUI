package canvas

// A butt cap stops exactly at the end point, a round one reaches half the
// width past it, and a square one reaches that same distance to the side as
// well — which is what keeps a dashed line the same length however it ends.
const (
	CapButt LineCap = iota
	CapRound
	CapSquare
)
