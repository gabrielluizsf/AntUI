package canvas

import "math"

// Matrix is a 2D affine transform in the CSS/SVG matrix(a, b, c, d, e, f)
// shape: (x', y') = (a*x + c*y + e, b*x + d*y + f). It is exactly what a
// transform: translate() rotate() scale() skew() matrix() list composes to,
// and the BlitMatrix path maps pixels through one.
type Matrix struct {
	A, B, C, D, E, F float64
}

// Identity is the transform that moves nothing.
func Identity() Matrix { return Matrix{A: 1, D: 1} }

// Translate shifts by tx to the right and ty down, the canvas's y growing down.
func Translate(tx, ty float64) Matrix { return Matrix{A: 1, D: 1, E: tx, F: ty} }

// Scale multiplies sizes by sx along x and sy along y. A negative factor
// mirrors the axis.
func Scale(sx, sy float64) Matrix { return Matrix{A: sx, D: sy} }

// Rotate spins clockwise by rad radians, as a y-down canvas reads it.
func Rotate(rad float64) Matrix {
	cos, sin := math.Cos(rad), math.Sin(rad)
	return Matrix{A: cos, B: sin, C: -sin, D: cos}
}

// Skew shears by the angles ax and ay, the CSS skew(ax, ay) shorthand for
// "skewX(ax) skewY(ay)" multiplied in that order. Angles near the vertical are
// pulled in, because tan(90°) is infinite.
func Skew(ax, ay float64) Matrix {
	x := Matrix{A: 1, C: math.Tan(clampTan(ax)), D: 1}
	y := Matrix{A: 1, B: math.Tan(clampTan(ay)), D: 1}
	return x.Mul(y)
}

// clampTan keeps a shear angle just off the vertical.
func clampTan(rad float64) float64 {
	const limit = math.Pi/2 - 1e-3
	if rad > limit {
		return limit
	}
	if rad < -limit {
		return -limit
	}
	return rad
}

// Mul returns the matrix product m·n: a point is transformed by n first and
// m afterwards, which is the order a CSS transform list multiplies in.
func (m Matrix) Mul(n Matrix) Matrix {
	return Matrix{
		A: m.A*n.A + m.C*n.B,
		B: m.B*n.A + m.D*n.B,
		C: m.A*n.C + m.C*n.D,
		D: m.B*n.C + m.D*n.D,
		E: m.A*n.E + m.C*n.F + m.E,
		F: m.B*n.E + m.D*n.F + m.F,
	}
}

// Map transforms one point.
func (m Matrix) Map(x, y float64) (float64, float64) {
	return m.A*x + m.C*y + m.E, m.B*x + m.D*y + m.F
}

// Inverse returns the transform that undoes m. A degenerate matrix — a zero
// determinant, as in scale(0) or a collapsed skew — has no inverse, and the
// bool reports it so a caller can paint nothing instead of dividing by zero.
func (m Matrix) Inverse() (Matrix, bool) {
	det := m.A*m.D - m.B*m.C
	if math.Abs(det) < 1e-12 {
		return Matrix{}, false
	}
	return Matrix{
		A: m.D / det,
		B: -m.B / det,
		C: -m.C / det,
		D: m.A / det,
		E: (m.C*m.F - m.D*m.E) / det,
		F: (m.B*m.E - m.A*m.F) / det,
	}, true
}
