package canvas

import "math"

// Rotate spins clockwise by rad radians, as a y-down canvas reads it.
func Rotate(rad float64) Matrix {
	cos, sin := math.Cos(rad), math.Sin(rad)
	return Matrix{A: cos, B: sin, C: -sin, D: cos}
}
