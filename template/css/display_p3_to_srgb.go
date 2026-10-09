package css

// displayP3ToSRGB maps a gamma-encoded display-p3 component triple into
// linear sRGB: linearise with the shared sRGB curve, cross the primaries into
// XYZ, and come back through the sRGB matrix.
func displayP3ToSRGB(r, g, b float64) (sr, sg, sb float64) {
	r, g, b = srgbToLinear(r), srgbToLinear(g), srgbToLinear(b)
	x := 0.4865709486482162*r + 0.2656676931690931*g + 0.1982172852343625*b
	y := 0.2289745640697488*r + 0.6917385218365064*g + 0.0792869140937449*b
	z := 0.0451133818589026*g + 1.0439443689009757*b
	return xyzToLinearSRGB(x, y, z)
}
