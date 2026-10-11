package svg

import "math"

// percentOf is what one axis of a percentage is a fraction of. In fractions of
// the shape every axis is the whole of one; in user units they are the drawing's
// own width and height, and an axis with no side to it — a radius — takes the
// diagonal.
func (out *gradient) percentOf(axis int) float64 {
	if !out.userSpace {
		return 1
	}
	switch axis {
	case 0:
		return out.viewport[0]
	case 1:
		return out.viewport[1]
	}
	return math.Hypot(out.viewport[0], out.viewport[1]) / math.Sqrt2
}
