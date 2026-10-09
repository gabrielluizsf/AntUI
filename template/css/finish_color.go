package css

import "github.com/gabrielluizsf/antui/canvas"

// finishColor turns CIE XYZ into a canvas colour through linear sRGB,
// honouring an optional alpha. Out-of-gamut components are clamped.
func finishColor(x, y, z float64, alpha string, hasAlpha bool) (canvas.Color, error) {
	r, g, b := xyzToLinearSRGB(x, y, z)
	return finishOK(linearToSRGB(r), linearToSRGB(g), linearToSRGB(b), alpha, hasAlpha)
}
