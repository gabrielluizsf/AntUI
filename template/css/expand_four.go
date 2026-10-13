package css

import "github.com/gabrielluizsf/antui/canvas"

// expandFour repeats a one-to-four CSS value across the four sides: one value
// fills all, two alternate top/bottom then left/right, three leave out the
// fourth (the left mirrors the right), four pass through.
//
// One function rather than one per element type — the rule is the same for a
// length, a border style and a colour, and only the values differ.
func expandFour[T any](v []T) [sideCount]T {
	switch len(v) {
	case 1:
		return [sideCount]T{v[0], v[0], v[0], v[0]}
	case 2:
		return [sideCount]T{v[0], v[1], v[0], v[1]}
	case 3:
		return [sideCount]T{v[0], v[1], v[2], v[1]}
	}
	return [sideCount]T{v[0], v[1], v[2], v[3]}
}

// expandFourColor is [expandFour] for the colours a box paints with, kept named
// because a colour reads better at a call site than a type argument does.
func expandFourColor(v []canvas.Color) [sideCount]canvas.Color { return expandFour(v) }
