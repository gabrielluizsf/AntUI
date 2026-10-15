package canvas

// applyFilter runs one filter on one colour, leaving its alpha alone.
func applyFilter(c Color, kind FilterKind, amount float64) Color {
	switch kind {
	case FilterGrayscale:
		return lerpMat(identityMat, grayscaleMat, clamp01(amount)).apply(c)
	case FilterSepia:
		return lerpMat(identityMat, sepiaMat, clamp01(amount)).apply(c)
	case FilterInvert:
		a := clamp01(amount)
		f := func(v uint8) uint8 { return clampByte(float64(v) + a*(255-2*float64(v))) }
		return RGBA(RGBAParams{R: f(c.R()), G: f(c.G()), B: f(c.B()), A: c.A()})
	case FilterBrightness:
		a := max(amount, 0)
		f := func(v uint8) uint8 { return clampByte(float64(v) * a) }
		return RGBA(RGBAParams{R: f(c.R()), G: f(c.G()), B: f(c.B()), A: c.A()})
	case FilterContrast:
		a := max(amount, 0)
		f := func(v uint8) uint8 { return clampByte((float64(v)-127.5)*a + 127.5) }
		return RGBA(RGBAParams{R: f(c.R()), G: f(c.G()), B: f(c.B()), A: c.A()})
	case FilterHueRotate:
		return hueRotateMat(amount).apply(c)
	}
	return c
}
