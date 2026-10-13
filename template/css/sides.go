package css

// The four sides, in the order every per-side array in a Style is written: the numbers
// the shorthand expands into, the index a border or a radius is stored at, and
// the order they go round the box. CSS writes it top, right, bottom, left —
// two values are top/bottom then left/right, three leave out the left because it
// mirrors the right — and that order is what [expandFour] hands back.
//
// The index is a plain int so it indexes the array directly. Nothing writes a
// side by number: the raw 0..3 are how these were read for years and they read
// as nothing at a call site.
const (
	SideTop    = 0
	SideRight  = 1
	SideBottom = 2
	SideLeft   = 3
)

// sideCount is how many sides a box has, which is how wide the arrays above are.
const sideCount = 4

// allSides names every side at once, for the shorthands that write the same
// thing on all four (border, margin, padding).
var allSides = [sideCount]int{SideTop, SideRight, SideBottom, SideLeft}

// oneSide names a single side, for the per-side shorthands (border-top and its
// three cousins) where a value lands on one side and leaves the rest alone. A
// -1 in the list is a side the shorthand did not name, which is what
// [parseBorderSides] reads to decide where a word goes.
func oneSide(side int) [sideCount]int {
	var s [sideCount]int
	for i := range s {
		s[i] = -1
	}
	s[side] = side
	return s
}
