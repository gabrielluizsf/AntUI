package canvas

// The two rules SVG and CSS ask for. NonZero is the default: a shape counts how
// many times it winds around a point and is inside wherever the count is not
// zero, so two nested rings drawn the same way are a ring and a hole. EvenOdd
// ignores the direction and counts the crossings, so any two nested rings leave
// a hole between them, however they were wound.
const (
	FillNonZero FillRule = iota
	FillEvenOdd
)
