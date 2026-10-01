package svg

import "github.com/gabrielluizsf/antui/canvas"

// Error is a drawing that could not be read at all, as opposed to one that was
// read and is only partly drawable — which is a list of warnings, not an error.
type Error struct{ What string }

// Error is the message a caller prints.
func (e *Error) Error() string { return "antui/svg: " + e.What }

// point is a canvas point with the two moves the smooth curve commands need: one
// that steps from where the line is, and one that mirrors the last control
// point through it, which is what `S` and `T` use for their first control.
type point = canvas.Point

// add is the point moved by another, used for the relative form of every
// command, where the numbers are measured from where the line is rather than
// from the origin.
func add(a, b point) point { return point{X: a.X + b.X, Y: a.Y + b.Y} }

// mirror is the point on the far side of a from b, which is the reflection of a
// about b and the first control point of a smooth curve that continues one.
func mirror(a, b point) point { return point{X: 2*b.X - a.X, Y: 2*b.Y - a.Y} }
