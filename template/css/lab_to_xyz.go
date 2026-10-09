package css

// The colour-space conversions below move a colour from a wide-gamut space
// into the canvas's sRGB. Reference whites are D65: X=0.95047, Y=1, Z=1.08883.
const (
	colorEpsilon = 216.0 / 24389.0
	colorKappa   = 24389.0 / 27.0
)

// labToXYZ maps CIE Lab into CIE XYZ under the D65 white point.
func labToXYZ(l, a, b float64) (x, y, z float64) {
	fy := (l + 16) / 116
	fx := fy + a/500
	fz := fy - b/200
	if v := fx * fx * fx; v > colorEpsilon {
		x = v
	} else {
		x = (116*fx - 16) / colorKappa
	}
	if v := fy * fy * fy; v > colorEpsilon {
		y = v
	} else {
		y = (116*fy - 16) / colorKappa
	}
	if v := fz * fz * fz; v > colorEpsilon {
		z = v
	} else {
		z = (116*fz - 16) / colorKappa
	}
	return x * 0.95047, y, z * 1.08883
}
