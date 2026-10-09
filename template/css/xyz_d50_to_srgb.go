package css

// xyzD50ToSRGB adapts a D50-adapted XYZ triple to D65 and maps it into
// linear sRGB, via the Bradford cone-response matrix.
func xyzD50ToSRGB(x, y, z float64) (r, g, b float64) {
	x = 0.9554734527*x - 0.0230985363*y + 0.0632593087*z
	y = -0.0283697069*x + 1.0099954581*y + 0.0210413984*z
	z = 0.0123140013*x - 0.0205076964*y + 1.3303659366*z
	return xyzToLinearSRGB(x, y, z)
}
