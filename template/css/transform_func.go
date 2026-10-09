package css

// TransformFunc is one function of a transform list. The zero value is a
// translate with no offset, and [Style.Transform] being nil means "none".
type TransformFunc struct {
	Kind uint8
	// translate() offsets; a percentage is of the border box.
	Dx, Dy Length
	// scale() factors.
	Sx, Sy float64
	// rotate() uses Ax; skew() uses both.
	Ax, Ay Angle
	// matrix() is the six numbers of the CSS/SVG matrix(a, b, c, d, e, f).
	M [6]float64
}
