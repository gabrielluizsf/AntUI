package canvas

import "math"

// colorMat is a 3x3 matrix applied to the red, green and blue channels of a
// pixel; alpha is never touched. It is the shape a CSS filter function takes.
type colorMat [9]float64

// identityMat leaves every channel where it is.
var identityMat = colorMat{1, 0, 0, 0, 1, 0, 0, 0, 1}

// grayscaleMat is the luminance matrix CSS grayscale() uses.
var grayscaleMat = colorMat{
	0.2126, 0.7152, 0.0722,
	0.2126, 0.7152, 0.0722,
	0.2126, 0.7152, 0.0722,
}

// sepiaMat is the matrix CSS sepia() uses.
var sepiaMat = colorMat{
	0.393, 0.769, 0.189,
	0.349, 0.686, 0.168,
	0.272, 0.534, 0.131,
}

// hueRotateMat is the SVG feColorMatrix hueRotate matrix for an angle in
// degrees.
func hueRotateMat(deg float64) colorMat {
	rad := deg * math.Pi / 180
	c, s := math.Cos(rad), math.Sin(rad)
	return colorMat{
		0.213 + c*0.787 - s*0.213, 0.715 - c*0.715 - s*0.715, 0.072 - c*0.072 + s*0.928,
		0.213 - c*0.213 + s*0.143, 0.715 + c*0.285 + s*0.140, 0.072 - c*0.072 - s*0.283,
		0.213 - c*0.213 - s*0.787, 0.715 - c*0.715 + s*0.715, 0.072 + c*0.928 + s*0.072,
	}
}

// apply runs the matrix on one colour.
func (m colorMat) apply(c Color) Color {
	r, g, b := float64(c.R()), float64(c.G()), float64(c.B())
	return RGBA(
		clampByte(m[0]*r+m[1]*g+m[2]*b),
		clampByte(m[3]*r+m[4]*g+m[5]*b),
		clampByte(m[6]*r+m[7]*g+m[8]*b),
		c.A(),
	)
}

// lerpMat mixes two matrices, t running from a to b.
func lerpMat(a, b colorMat, t float64) colorMat {
	var out colorMat
	for i := range out {
		out[i] = a[i] + (b[i]-a[i])*t
	}
	return out
}
