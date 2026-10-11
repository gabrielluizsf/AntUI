package canvas

import "math"

// Skew shears by the angles ax and ay, the CSS skew(ax, ay) shorthand for
// "skewX(ax) skewY(ay)" multiplied in that order. Angles near the vertical are
// pulled in, because tan(90°) is infinite.
func Skew(ax, ay float64) Matrix {
	x := Matrix{A: 1, C: math.Tan(clampTan(ax)), D: 1}
	y := Matrix{A: 1, B: math.Tan(clampTan(ay)), D: 1}
	return x.Mul(y)
}
