package canvas

// boxBlur runs one moving-average pass, horizontally or vertically, clamping
// at the edges so the picture does not darken as it approaches them. The result
// goes into dst, which the caller owns and hands back on the next pass.
func boxBlur(src, dst []float64, w, h, r int, horizontal bool) []float64 {
	span := float64(2*r + 1)
	if horizontal {
		for row := range h {
			boxBlurLine(src, dst, row*w*4, 4, w, r, span)
		}
		return dst
	}
	for col := range w {
		boxBlurLine(src, dst, col*4, w*4, h, r, span)
	}
	return dst
}
