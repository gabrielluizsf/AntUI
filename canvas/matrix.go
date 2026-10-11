package canvas

// Matrix is a 2D affine transform in the CSS/SVG matrix(a, b, c, d, e, f)
// shape: (x', y') = (a*x + c*y + e, b*x + d*y + f). It is exactly what a
// transform: translate() rotate() scale() skew() matrix() list composes to,
// and the BlitMatrix path maps pixels through one.
type Matrix struct {
	A, B, C, D, E, F float64
}
