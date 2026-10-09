package css

// xyzToLinearSRGB maps CIE XYZ into linear sRGB; negative or out-of-range
// components are clamped by the caller's byte conversion.
func xyzToLinearSRGB(x, y, z float64) (r, g, b float64) {
	return 3.2404542*x - 1.5371385*y - 0.4985314*z,
		-0.9692660*x + 1.8760108*y + 0.0415560*z,
		0.0556434*x - 0.2040259*y + 1.0572252*z
}
